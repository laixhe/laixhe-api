<?php

namespace App\Http\Services;

use Throwable;
use RuntimeException;

use Godruoyi\Snowflake\Snowflake;
use Illuminate\Database\QueryException;
use Illuminate\Support\Facades\DB;

use App\Http\Requests\LoginRequest;
use App\Http\Requests\RegisterRequest;
use App\Result\ResultCode;
use App\Models\User;
use App\Models\UserExtend;
use App\Models\UserThirdParty;

/**
 * 鉴权服务相关 (与 Go 端 services/auth.go 对齐)
 */
class AuthService
{
    /**
     * 注册 (与 Go 端 Auth.Register 对齐)
     *
     * 先查邮箱是否已注册, 避免无效的 bcrypt 计算;
     * email 为唯一索引 (与 webapi.sql 一致), 先查后插 + 数据库唯一约束双重防重。
     * 在同一事务中创建用户、扩展信息、第三方关联。
     *
     * @param RegisterRequest $req
     * @return array
     *
     * @throws RuntimeException
     */
    public function register(RegisterRequest $req): array
    {
        if (User::query()->where('email', $req->email)->exists()) {
            throw new RuntimeException('邮箱已存在', ResultCode::Param->value);
        }
        try {
            $user = DB::transaction(function () use ($req): User {
                $password = password_hash($req->password, PASSWORD_BCRYPT);
                if ($password === false) {
                    // bcrypt 计算失败 (极少数场景), 与 Go 端 hash 失败处理对齐
                    throw new RuntimeException('', ResultCode::Service->value);
                }
                $user = User::query()->create([
                    'type_id' => 1, // 普通用户
                    'account' => (string)app(Snowflake::class)->id(), // 全局唯一账号 (复用容器单例)
                    'mobile' => '',
                    'email' => $req->email,
                    'password' => $password,
                    'nickname' => $req->nickname,
                    'avatar_url' => '',
                    'sex' => 0, // 未填写
                    'states' => 1, // 正常
                ]);
                UserExtend::query()->create(['uid' => $user->id]);
                UserThirdParty::query()->create(['uid' => $user->id]);
                return $user;
            });
            return $user->toArray();
        } catch (QueryException $e) {
            // 唯一键冲突 (account/email/关联表 uid): 并发注册同邮箱等极端情况才触发
            if ($e->getCode() === '23000') {
                throw new RuntimeException('注册失败，请稍后再试', ResultCode::Param->value);
            }
            throw new RuntimeException($e->getMessage(), ResultCode::Service->value);
        } catch (Throwable $e) {
            throw new RuntimeException($e->getMessage(), ResultCode::Service->value);
        }
    }

    // 时序侧信道防护用的假 bcrypt 哈希 (每次进程首次登录时惰性生成并缓存)
    private static ?string $dummyHash = null;

    private static function dummyHash(): string
    {
        return self::$dummyHash ??= password_hash('dummy-password', PASSWORD_BCRYPT);
    }

    /**
     * 登录 (与 Go 端 Auth.Login 对齐)
     *
     * 封禁账号与密码错误返回同一提示, 避免暴露账号状态 (可被探测)。
     * 账号不存在 / 被封禁时也执行一次 password_verify (对假哈希):
     * 使各失败路径的响应时间一致, 防止通过时序差异枚举账号
     * (文案相同只能防文本探测, 不能防时序侧信道)。
     *
     * @param LoginRequest $req
     * @return array 未找到账号或账号被封禁时返回空数组
     */
    public function login(LoginRequest $req): array
    {
        // select * from `user` where `email` = ? limit 1
        $user = User::query()->where('email', $req->email)->first();
        if (empty($user)) {
            // 账号不存在: 执行一次 bcrypt 校验使耗时与"密码错误"路径一致
            password_verify($req->password, self::dummyHash());
            return [];
        }
        // 账号状态: 1=正常 0=封禁 (见迁移文件 comment); 封禁账号与密码错误返回同一提示
        if ((int)$user->states !== 1) {
            // 封禁账号: 同样执行一次 bcrypt 校验保持时序一致
            password_verify($req->password, self::dummyHash());
            return [];
        }
        // 登录校验需要密码哈希: 临时解除 password 的隐藏, 供控制器 password_verify 使用
        return $user->makeVisible('password')->toArray();
    }
}
