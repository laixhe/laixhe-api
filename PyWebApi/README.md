### 创建虚拟环境(确保当前目录有 pyproject.toml 文件)
```
poetry install
```

### 启动服务
 - 启动命令中的 app 指的是 app = FastAPI() 变量，也可以是其他自己定义的名称
 - 
```
uvicorn main:app --reload

http://127.0.0.1/docs
http://127.0.0.1/redoc
http://127.0.0.1/openapi.json
```

```
├── README.md  # 项目介绍
├── app
│   ├── __init__.py
│   ├── config  # 配置相关
│   │   └── __init__.py
│   ├── constant  # 常量相关
│   │   └── __init__.py
│   ├── controller  # 常量相关
│   │   └── __init__.py
│   ├── dependencies  # 封装被依赖函数
│   │   └── __init__.py
│   ├── middleware # 中间件
│   │   └── __init__.py
│   ├── models # 数据模型文件，和表结构对应
│   │   └── __init__.py
│   ├── service # 就具体业务实现逻辑
│   │   └── __init__.py
│   └── utils # 工具类
│       ├── __init__.py
├── main.py # 主文件
```

### Pydantic 
- 基本数据类型: int, float, str, bool
- 可选参数: Optional[type] 表示可选参数,Union[x, None]也可以表示可选
- 整数范围: 结合 conint 函数判断数字范围 ,如 age: conint(ge=18, le=30); ge:大于等于、gt:大于、le:小于等于、lt:小于
- 字符长度: 结合 constr 函数判断字符长度,如: constr(min_length=6, max_length=10)
- 正则表达式: 使用 constr 函数中的参数 regex ，可以用于进行正则表达式验证
- 枚举验证: 使用 Enum 定义枚举类,验证
- 列表类型: 使用 List[type] 来限制列表值的类型，并尝试把参数转成对应的类型
- 字典类型: Dict[key_type, value_type] 来限制字典 key 和 val 类型，并尝试把参数转成对应的类型
- 其他验证
  - EmailStr: 用于验证字符串是否是有效的电子邮件地址。
  - IPvAnyAddress: 用于验证字符串是否是有效的 IPv4 或 IPv6 地址。
  - StrictBool： 用于验证字符串是否是严格的布尔值（true 或 false）。
  - AnyHttpUrl: 用于验证字符串是否是有效的 URL，包括以 http 或 https 开头的URL

```
from enum import Enum
from typing import Union, Optional, List, Dict

# 导入 pydantic 对应的模型基类
from pydantic import BaseModel, constr, conint

class GenderEnum(str, Enum):
    """
    性别枚举
    """
    male = "男"
    female = "女"

class PydanticVerifyParam(BaseModel):
    """
    用来学习使用 pydantic 模型验证
    """
    user_name: str  # 基本类型
    age: conint(ge=18, le=30)  # 整数范围：18 <= age <= 30
    password: constr(min_length=6, max_length=10)  # 字符长度
    phone: constr(regex=r'^1\d{10}$')  # 正则验证手机号
    address: Optional[str] = None  # 可选参数
    sex: GenderEnum  # 枚举验证,只能传: 男和女
    likes: List[str]  # 值会自动转成传字符串列表
    scores: Dict[str, float]  # key 会转成字符串 val 会转成浮点型
```