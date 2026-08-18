<?php

namespace App\Models;

/**
 * 通用配置表
 */
class ConfigCommon extends BaseModel
{
    // 与模型关联的数据表名
    protected $table = 'config_common';
    // 与数据表关联的主键
    protected $primaryKey = 'id';
    // 指明模型的ID是否自动递增
    public $incrementing = true;
    // 自动递增ID的数据类型
    protected $keyType = 'integer';
    // config_common 无 created_at/updated_at 字段
    public $timestamps = false;

    // 可修改的表字段 (批量赋值白名单)
    protected $fillable = [
        'key',
        'value',
        'describe',
    ];
}
