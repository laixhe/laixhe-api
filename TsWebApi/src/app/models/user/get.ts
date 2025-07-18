import type { User } from "@prisma/client";

import prisma from "@core/prisma";

export default async function GetByEmail(email: string): Promise<User | null> {
  return await prisma.user.findUnique({
    where: {
      email: email,
    },
  });
}
