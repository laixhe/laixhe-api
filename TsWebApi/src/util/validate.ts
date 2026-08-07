// 邮箱格式正则
const emailRe = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

// 邮箱格式校验
export function isEmail(email: string): boolean {
  return emailRe.test(email);
}

// 密码长度校验（>= 6 位）
export function isPasswordTooShort(password: string): boolean {
  return password.length < 6;
}

// 密码字符正则：仅允许字母、数字、_、@、$，长度 >= 6
const passwordRe = /^[a-zA-Z0-9_@$]{6,}$/;

// 密码字符校验（仅允许字母、数字、_、@、$）
export function isPasswordInvalid(pattern: string): boolean {
  return !passwordRe.test(pattern);
}

// 昵称过短校验（< 2 字）
export function isNicknameTooShort(nickname: string): boolean {
  return nickname.length < 2;
}

// 昵称过长校验（> 20 字）
export function isNicknameTooLong(nickname: string): boolean {
  return nickname.length > 20;
}
