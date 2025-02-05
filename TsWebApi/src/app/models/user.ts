import type { User } from "@prisma/client";

import prisma from "@core/prisma";

async function GetByEmail(email: string): Promise<User | null> {
  return await prisma.user.findUnique({
    where: {
      email: email,
    },
  });
}

export const UserModel = {
  GetByEmail,
};
