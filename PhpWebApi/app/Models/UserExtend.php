<?php

namespace App\Models;

/**
 * 用户扩展表
 */
class UserExtend extends BaseModel
{
    // 与模型关联的数据表名
    protected $table = 'user_extend';
    // 与数据表关联的主键
    protected $primaryKey = 'id';
    // 指明模型的ID是否自动递增
    public $incrementing = true;
    // 自动递增ID的数据类型
    protected $keyType = 'integer';
    // 指示模型是否主动维护时间戳
    public $timestamps = false;

    // 可修改的表字段
    protected $fillable = [
        'uid',
        'birthday',
        'height',
        'weight',
    ];
}
