// 用户类型
export enum UserType {
  Ordinary = 1, // 普通
}

// 校验用户类型
export function isUserType(value: number): value is UserType {
  return Object.values(UserType).includes(value);
}

export function getUserTypeText(type: UserType): string {
  switch (type) {
    case UserType.Ordinary:
      return "普通";
    default:
      return "";
  }
}

// 用户状态
export enum UserStates {
  Banned = 0, // 封禁
  Normal = 1, // 正常
}

// 校验用户状态
export function isUserStates(value: number): value is UserStates {
  return Object.values(UserStates).includes(value);
}

export function getUserStatesText(state: UserStates): string {
  switch (state) {
    case UserStates.Banned:
      return "封禁";
    case UserStates.Normal:
      return "正常";
    default:
      return "";
  }
}

// 用户性别
export enum UserSex {
  Male = 1, // 男
  Female = 2, // 女
}

// 校验用户性别
export function isUserSex(value: number): value is UserSex {
  return Object.values(UserSex).includes(value);
}

export function getUserSexText(sex: UserSex): string {
  switch (sex) {
    case UserSex.Male:
      return "男";
    case UserSex.Female:
      return "女";
    default:
      return "";
  }
}

