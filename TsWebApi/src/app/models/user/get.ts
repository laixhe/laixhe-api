import type { User } from "@generated/prisma/client";
import { prisma } from "@core/prisma";

export async function getByUid(uid: number): Promise<User | null | undefined> {
  return await prisma.user.findUnique({
    where: {
      id: uid,
    },
  });
}

export async function getByMobile(
  mobile: string
): Promise<User | null | undefined> {
  return await prisma.user.findFirst({
    where: {
      mobile: mobile,
    },
  });
}

export async function getByEmail(
  email: string
): Promise<User | null | undefined> {
  return await prisma.user.findFirst({
    where: {
      email: email,
    },
  });
}
