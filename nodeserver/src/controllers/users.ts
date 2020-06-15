import { Request, Response } from "express";
import { bankService, Role, userService } from "../models/bank";
import { authService } from "./auth";

export const list = (req: Request, res: Response) => {
  if (!authService.allowed(req, "users_list")) {
    res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
    return;
  }
  const users = userService.userList();
  res.status(200).json({ users });
};

export const create = (req: Request, res: Response) => {
  if (!authService.allowed(req, "users_create")) {
    res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
    return;
  }

  const fullname = req.body.fullname;
  const username = req.body.username;
  const password = req.body.password;

  const user = bankService.findUserByUserName(username);
  if (user !== null) {
    res
      .status(500)
      .json({ errorcode: 500, errormessage: "User Already Exists" });
    return;
  }

  const userid = userService.create(fullname, username, password);

  res.status(200).json({ userid });
};

export const get = (req: Request, res: Response) => {
  const userid = req.params.id;

  if (!authService.allowed(req, "users_get")) {
    res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
    return;
  }
  if (req.userRole === Role.User && req.userID !== userid) {
    res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
    return;
  }
  const user = bankService.findUserByID(userid);
  if (user === null) {
    res.status(500).json({ errorcode: 500, errormessage: "No such user" });
    return;
  }

  res.status(200).json({ user });
};
