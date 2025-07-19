// 用户类型
enum UserType {
  ordinary = 1, // 普通
}

// 校验用户类型
function isUserType(value: number): value is UserType {
  return Object.values(UserType).includes(value);
}

// 用户状态
enum UserStates {
  Banned = 0, // 封禁
  Normal = 1, // 正常
}

// 校验用户状态
function isUserStates(value: number): value is UserStates {
  return Object.values(UserStates).includes(value);
}

// 用户性别
enum UserSex {
  Male = 1, // 男
  Female = 2, // 女
}

// 校验用户性别
function isUserSex(value: number): value is UserSex {
  return Object.values(UserSex).includes(value);
}

export { UserType, UserStates, UserSex, isUserType, isUserStates, isUserSex };
