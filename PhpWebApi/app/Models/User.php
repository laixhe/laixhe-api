<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Notifications\Notifiable;

/**
 * 用户表
 */
class User extends BaseModel
{
    use HasFactory;
    use Notifiable;

    // 与模型关联的数据表名
    protected $table = 'user';
    // 与数据表关联的主键
    protected $primaryKey = 'id';
    // 指明模型的ID是否自动递增
    public $incrementing = true;
    // 自动递增ID的数据类型
    protected $keyType = 'integer';
    // 指示模型是否主动维护时间戳 (需要 created_at 和 updated_at 字段存在你的模型数据表中)
    public $timestamps = true;

    // 可修改的表字段 (批量赋值白名单)
    protected $fillable = [
        'type_id',
        'account',
        'mobile',
        'email',
        'password',
        'nickname',
        'avatar_url',
        'sex',
        'states',
        self::UPDATED_AT,
    ];

    // 序列化为数组/JSON 时隐藏敏感字段 (password 仍可通过 $user->password 访问, 用于登录校验)
    protected $hidden = ['password'];

    /**
     * 查询用户时使用的列 (排除 password, 与 Go 端 UserColumnsNoPassword 对齐)。
     * 密码哈希只在登录校验时需要, 其它场景读取会把无用的敏感字段拉进内存。
     *
     * @return string[]
     */
    public static function noPassword(): array
    {
        return [
            'id',
            'type_id',
            'account',
            'mobile',
            'email',
            'nickname',
            'avatar_url',
            'sex',
            'states',
            'created_at',
            'updated_at',
        ];
    }

}
