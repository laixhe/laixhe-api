package core

import (
	"strings"
	"sync"
	"testing"
)

// TestBcryptPoolHashAndCheck 验证 worker 池的 Hash/Check 与 gonet/crypto 行为一致:
// 同一密码两次哈希不同 (随机盐), 但均可通过 Check 校验
func TestBcryptPoolHashAndCheck(t *testing.T) {
	p := NewBcryptPool(2)
	defer p.Close()

	hash1, err := p.Hash("mypassword123")
	if err != nil {
		t.Fatalf("Hash 失败: %v", err)
	}
	if hash1 == "" {
		t.Fatal("Hash 不应为空")
	}
	if !strings.HasPrefix(hash1, "$2") {
		t.Fatalf("哈希应以 $2 开头, got: %s", hash1)
	}
	hash2, err := p.Hash("mypassword123")
	if err != nil {
		t.Fatalf("第二次 Hash 失败: %v", err)
	}
	if hash1 == hash2 {
		t.Fatal("同一密码两次哈希应不同 (随机盐)")
	}

	if !p.Check("mypassword123", hash1) {
		t.Fatal("正确密码应通过 Check")
	}
	if !p.Check("mypassword123", hash2) {
		t.Fatal("第二份哈希同样应通过 Check")
	}
	if p.Check("wrongpassword", hash1) {
		t.Fatal("错误密码不应通过 Check")
	}
}

// TestBcryptPoolEmptyPassword 验证空密码边界: 空密码可正常哈希, 空哈希 Check 返回 false
func TestBcryptPoolEmptyPassword(t *testing.T) {
	p := NewBcryptPool(1)
	defer p.Close()

	hash, err := p.Hash("")
	if err != nil {
		t.Fatalf("空密码哈希失败: %v", err)
	}
	if !p.Check("", hash) {
		t.Fatal("空密码与自身哈希应匹配")
	}
	if p.Check("any", "") {
		t.Fatal("空哈希不应匹配任何密码")
	}
}

// TestBcryptPoolConcurrent 验证高并发提交下结果不串号:
// 每个 goroutine 用自己的唯一密码, 校验必须只对自己成立
func TestBcryptPoolConcurrent(t *testing.T) {
	p := NewBcryptPool(4)
	defer p.Close()

	const n = 32
	var wg sync.WaitGroup
	errCh := make(chan string, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(seq int) {
			defer wg.Done()
			password := strings.Repeat("p", 6+seq%10) + "-" + string(rune('a'+seq))
			hash, err := p.Hash(password)
			if err != nil {
				errCh <- "hash err: " + err.Error()
				return
			}
			if !p.Check(password, hash) {
				errCh <- "自身密码未通过校验"
				return
			}
			if p.Check(password+"x", hash) {
				errCh <- "错误密码不应通过校验"
			}
		}(i)
	}
	wg.Wait()
	close(errCh)
	for msg := range errCh {
		t.Fatal(msg)
	}
}
