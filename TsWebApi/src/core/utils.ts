// 手机号格式
const MobileRegex = /^1[2-9]\d{9}$/;
// 身份证格式
const IdentityNumberRegex =
  /^[1-9]\d{5}(18|19|20)\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$/;

// 验证手机号格式
export function validateMobile(mobile: string): boolean {
  if (mobile.length != 11) {
    return false;
  }
  return MobileRegex.test(mobile);
}

// 验证身份证格式
export function validateIdentityNumber(identityNumber: string): boolean {
  if (identityNumber.length != 18) {
    return false;
  }
  return IdentityNumberRegex.test(identityNumber);
}

// 从身份证号码中提取生日
export function extractBirthdayFromIdentityNumber(identityNumber: string): string {
  if (!validateIdentityNumber(identityNumber)) {
    return "";
  }
  return identityNumber.substring(6, 14);
}
