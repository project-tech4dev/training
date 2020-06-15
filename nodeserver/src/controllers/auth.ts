import crypto from "crypto";
import { Request, Response } from "express";
import { bankService, Role, User } from "../models/bank";

interface AuthService {
  isAuth: (a: string, b: string) => boolean;
  getToken: (a: string, b: string) => string;
  createSession: (a: string, b: User) => void;
  sessionExists: (a: string) => boolean;
  getSessionUser: (a: string) => User;
  allowed: (a: Request, b: string, c?: string) => boolean;
}

interface Session {
  [index: string]: User;
}

export class ASClass implements AuthService {
  session: Session = {};
  allowed = (req: Request, permission: string): boolean => {
    let ret = false;
    switch (permission) {
      case "accounts_list":
      case "accounts_create":
      case "users_list":
        ret = req.userRole === Role.BankManager;
        break;
      case "accounts_credit":
      case "accounts_debit":
        ret = req.userRole !== Role.BankManager;
        break;
      default:
        ret = true;
        break;
    }
    return ret;
  };

  isAuth = (username: string, password: string): boolean => {
    let repl = true;
    const u = bankService.findUserByUserName(username);

    if (u) {
      if (u.password !== password) {
        repl = false;
      }
    } else {
      repl = false;
    }
    return repl;
  };

  getToken = (username: string, password: string): string => {
    const s = crypto
      .createHash("sha256")
      .update(password)
      .digest("hex") as string;
    return s;
  };

  createSession = (token: string, user: User) => {
    this.session[token] = user;
  };

  sessionExists = (token: string): boolean => {
    return !!this.session[token];
  };

  getSessionUser = (token: string): User => {
    return this.session[token];
  };
}

export let authService = new ASClass();

export const login = (req: Request, res: Response) => {
  const username = req.body.username;
  const password = req.body.password;

  if (authService.isAuth(username, password)) {
    const token = authService.getToken(username, password);
    const user = bankService.findUserByUserName(username);
    authService.createSession(token, user);
    res.status(200).json({ id: user.userid, token });
  } else {
    res.status(500).send("Username or Password don't match");
  }
};
