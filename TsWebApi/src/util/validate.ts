// 邮箱格式正则
const emailRe = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

export function isEmail(email: string): boolean {
  return emailRe.test(email);
}

// 密码长度校验（>= 6 位）
export function isPasswordTooShort(password: string): boolean {
  return password.length < 6;
}

// 密码字符正则：仅允许字母、数字、_、@、$，长度 >= 6
const passwordRe = /^[a-zA-Z0-9_@$]{6,}$/;

// 说明：下面两个密码校验函数在"长度 < 6"这一规则上重叠
// （passwordRe 自带 {6,} 长度限制），分开定义只是为了：
//   - isPasswordTooShort：单独命中"过短"这一种情况，给出专属提示
//   - isPasswordInvalid：命中字符集非法（含过短），提示覆盖完整规则
// 调用方先查前者、再查后者，即可返回精确的错误提示。

// 密码字符校验（仅允许字母、数字、_、@、$，且长度 >= 6）
export function isPasswordInvalid(password: string): boolean {
  return !passwordRe.test(password);
}

// 昵称过短校验（< 2 字）
export function isNicknameTooShort(nickname: string): boolean {
  return nickname.length < 2;
}

// 昵称过长校验（> 20 字）
export function isNicknameTooLong(nickname: string): boolean {
  return nickname.length > 20;
}
