package core

import (
	"runtime"
	"sync"

	"github.com/laixhe/gonet/crypto"
)

// bcryptTask 单个 bcrypt 计算任务 (由请求方提交, worker 完成后写回结果)
type bcryptTask struct {
	isCheck  bool   // true=校验密码, false=生成哈希
	password string // 明文密码 (Hash/Check 均需要)
	hash     string // 仅 Check 使用: 待校验的 bcrypt 哈希
	result   chan bcryptResult
}

// bcryptResult bcrypt 任务结果
type bcryptResult struct {
	hash string // 仅 Hash 使用: 生成的哈希
	ok   bool   // 仅 Check 使用: 密码是否匹配
	err  error  // 仅 Hash 使用: 计算失败原因
}

// BcryptPool bcrypt 计算 worker 池
//
// 背景: bcrypt cost=10 为 CPU 密集计算 (单次约 50-100ms), 若直接在请求 goroutine 上
// 同步执行, 高并发注册/登录会占满 GOMAXPROCS 个线程, 拖垮同进程内的其它请求。
// 此实现参照 Rust 版 spawn_blocking 思路: 固定 N 个 worker goroutine (默认 = GOMAXPROCS),
// 请求方将任务提交到 channel 后阻塞等待结果, worker 完成计算后写回。
//
// 与 Rust spawn_blocking (任务随到随跑、由运行时调度) 的差异: 池大小固定,
// 并发任务超过池大小时请求会排队等待, 属于"背压"而非"无限并行",
// 同时配合接口限流 (默认 1000 次/分钟/IP, 见 config.yaml 的 limit 配置) 防止滥用。
//
// 注意: Hash/Check 的输入输出均为值拷贝, 池设计为进程生命周期内常驻单例,
// 无需每次请求创建; 见 Server.Bcrypt。
type BcryptPool struct {
	// tasks 任务提交 channel (带缓冲: 缓冲为池大小的 2 倍, 减少请求方提交时的锁竞争)
	tasks chan bcryptTask
	// resultPool 结果 channel 复用池: Hash/Check 每调用一次都新建 channel 会产生堆分配,
	// 用 sync.Pool 缓存已回收的 channel, 减少高频注册/登录路径的 GC 压力。
	// 归还时 channel 保证为空 (缓冲 1, 结果已被接收方读出), 可安全复用。
	resultPool sync.Pool
}

// NewBcryptPool 创建 bcrypt worker 池
//
// workers 指定 worker goroutine 数量, 由 config.yaml 的 bcrypt.workers 配置 (见 Config.Bcrypt):
// - 传 0 或负数: 自动取 runtime.GOMAXPROCS(0) (当前 CPU 逻辑核数);
// - 显式配置: 按部署机器规格设置, 建议不超过核数。
//
// 生产调优: bcrypt 为纯 CPU 计算, 池大于核数无吞吐收益反而增加线程切换开销;
// 池小于核数时, 高并发注册/登录会在提交处排队 (背压), 表现为接口延迟上升。
// 可按实际压测调整: 观察 CPU 使用率与注册/登录接口 P99 延迟,
// 若 CPU 未打满但接口排队延迟明显, 可适当增大 workers; 反之调小。
func NewBcryptPool(workers int) *BcryptPool {
	if workers <= 0 {
		workers = runtime.GOMAXPROCS(0)
	}
	p := &BcryptPool{
		tasks: make(chan bcryptTask, workers*2),
		resultPool: sync.Pool{
			New: func() any {
				return make(chan bcryptResult, 1)
			},
		},
	}
	for i := 0; i < workers; i++ {
		go p.worker()
	}
	return p
}

// worker 单个 worker goroutine: 循环取任务执行 bcrypt, 结果写回 result channel
func (p *BcryptPool) worker() {
	for task := range p.tasks {
		if task.isCheck {
			// bcrypt verify 与 hash 同为 CPU 密集, 放 worker 池执行 (与 Hash 对称)
			task.result <- bcryptResult{ok: crypto.BcryptPasswordCheck(task.password, task.hash)}
		} else {
			hash, err := crypto.BcryptPasswordHash(task.password)
			task.result <- bcryptResult{hash: hash, err: err}
		}
	}
}

// Hash 在 worker 池上计算密码的 bcrypt 哈希 (阻塞直到结果返回)
func (p *BcryptPool) Hash(password string) (string, error) {
	result := p.resultPool.Get().(chan bcryptResult)
	task := bcryptTask{
		password: password,
		result:   result, // 缓冲 1: worker 写入后无需等待接收方即可退出
	}
	p.tasks <- task
	res := <-task.result
	// 结果已读出, channel 为空, 归还池中复用 (见 resultPool 注释)
	p.resultPool.Put(result)
	return res.hash, res.err
}

// Check 在 worker 池上校验密码与 bcrypt 哈希是否匹配 (阻塞直到结果返回)
func (p *BcryptPool) Check(password, hash string) bool {
	result := p.resultPool.Get().(chan bcryptResult)
	task := bcryptTask{
		isCheck:  true,
		password: password,
		hash:     hash,
		result:   result, // 缓冲 1: worker 写入后无需等待接收方即可退出
	}
	p.tasks <- task
	res := <-task.result
	p.resultPool.Put(result)
	return res.ok
}

// Close 关闭 worker 池 (等待所有 worker 退出)
//
// 生产环境由 NewServer 创建单例并随进程存活, 无需调用;
// 主要用于测试场景, 避免 goroutine 泄漏。
func (p *BcryptPool) Close() {
	close(p.tasks)
}
