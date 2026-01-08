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

//    // 可修改的表字段
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

//    // 隐藏查询结果的字段
//    protected $hidden = ['password'];

    public function getTable(): string
    {
        return $this->table;
    }

}
