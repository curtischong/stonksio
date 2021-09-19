import apiClient from "./client";
import { Balance } from "../types";

export const getWallet = async (username: string): Promise<Balance> => {
  const response = await apiClient().get("/wallet", {
    params: {
      username,
    },
  });
  return response.data;
};
