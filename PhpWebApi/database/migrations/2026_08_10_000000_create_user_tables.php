<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

/**
 * 核心业务表 (与仓库根目录 webapi.sql / Rust 端 docs/schema.sql 保持一致):
 * - user: 用户主表, account / email 唯一, 注册先查后插 + 数据库唯一约束双重防重
 * - user_extend / user_third_party: 用户扩展与第三方关联, uid 唯一
 * - config_common: 通用配置
 */
return new class extends Migration
{
    public function up(): void
    {
        Schema::create('user', function (Blueprint $table) {
            $table->unsignedInteger('id')->autoIncrement();
            $table->integer('type_id')->default(0)->comment('类型 1普通');
            $table->string('account', 120)->default('')->comment('账号');
            $table->string('mobile', 120)->default('')->comment('手机号');
            $table->string('email', 120)->default('')->comment('邮箱');
            $table->string('password', 120)->default('')->comment('密码');
            $table->string('nickname', 120)->default('')->comment('昵称');
            $table->string('avatar_url', 255)->default('')->comment('头像地址');
            $table->integer('sex')->default(0)->comment('性别 0未填写 1男 2女');
            $table->integer('states')->default(0)->comment('状态 0封禁 1正常');
            $table->dateTime('created_at')->comment('创建时间');
            $table->dateTime('updated_at')->comment('更新时间');

            // 唯一索引: 账号全局唯一 (对齐 webapi.sql / Go / Rust / TS 端)
            $table->unique('account', 'user_account_idx');
            $table->index('mobile', 'user_mobile_idx');
            // email 唯一索引 (与 webapi.sql 一致; 注册先查后插 + 数据库唯一约束双重防重)
            $table->unique('email', 'user_email_idx');
            $table->comment('用户');
        });

        Schema::create('user_extend', function (Blueprint $table) {
            $table->unsignedInteger('id')->autoIncrement();
            $table->unsignedInteger('uid')->comment('用户ID');
            $table->integer('birthday')->default(0)->comment('生日(年月日)');
            $table->integer('height')->default(0)->comment('身高(cm)');
            $table->integer('weight')->default(0)->comment('体重(kg)');

            // 一对一约束: 一个用户一条扩展记录 (对齐 webapi.sql / Go / Rust / TS 端)
            $table->unique('uid', 'user_extend_uid_idx');
            $table->comment('用户扩展');
        });

        Schema::create('user_third_party', function (Blueprint $table) {
            $table->unsignedInteger('id')->autoIncrement();
            $table->unsignedInteger('uid')->comment('用户ID');
            $table->string('wechat_unionid', 200)->default('')->comment('微信unionid');
            $table->string('wechat_openid', 200)->default('')->comment('微信openid');

            // 一对一约束: 一个用户一条第三方绑定记录 (对齐 webapi.sql / Go / Rust / TS 端)
            $table->unique('uid', 'user_third_party_uid_idx');
            $table->index('wechat_openid', 'user_third_party_wechat_openid_idx');
            $table->comment('用户第三方');
        });

        Schema::create('config_common', function (Blueprint $table) {
            $table->unsignedInteger('id')->autoIncrement();
            $table->string('key', 255)->default('')->comment('配置键');
            $table->string('value', 512)->default('')->comment('配置值');
            $table->string('describe', 255)->default('')->comment('描述');

            $table->index('key', 'config_common_key_idx');
            $table->comment('通用配置');
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('config_common');
        Schema::dropIfExists('user_third_party');
        Schema::dropIfExists('user_extend');
        Schema::dropIfExists('user');
    }
};
