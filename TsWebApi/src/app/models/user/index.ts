import {
  UserType,
  UserStates,
  UserSex,
  isUserType,
  isUserStates,
  isUserSex,
  getUserTypeText,
  getUserStatesText,
  getUserSexText,
} from "@model/user/user";
import { getByUid, getByMobile, getByEmail } from "@model/user/get";

export const UserModel = {
  UserType,
  UserStates,
  UserSex,
  isUserType,
  isUserStates,
  isUserSex,
  getUserTypeText,
  getUserStatesText,
  getUserSexText,
  getByUid,
  getByMobile,
  getByEmail,
};
