import { Request, Response } from "express";
import { bankService, Role } from "../models/bank";
import { authService } from "./auth";

export const list = (req: Request, res: Response) => {
  if (!authService.allowed(req, "accounts_list")) {
    res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
    return;
  }

  const userlist = bankService.accountList();
  res.status(200).json({ accounts: userlist });
};

export const create = (req: Request, res: Response) => {
  if (!authService.allowed(req, "accounts_create")) {
    res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
    return;
  }

  const userid = req.body.userid as string;
  const balance = req.body.balance as number;

  const u = bankService.findUserByID(userid);
  if (u === null) {
    res.status(500).json({ errorcode: 500, errormessage: "No such user" });
    return;
  }
  const account = bankService.createAccount(userid, balance);

  res.status(200).json({ accountid: account.id, balance: account.balance });
};

export const get = (req: Request, res: Response) => {
  const acctid = req.params.id;

  if (!authService.allowed(req, "accounts_get")) {
    res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
    return;
  }
  const acct = bankService.findAccountByID(acctid);
  if (acct === null) {
    res.status(500).json({ errorcode: 500, errormessage: "No such account" });
    return;
  }
  if (req.userRole === Role.User && req.userID !== acct.userid) {
    res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
    return;
  }

  res.status(200).json({
    accountid: acct.id,
    balance: acct.balance,
    activity: acct.activity,
  });
};

export const credit = (req: Request, res: Response) => {
  const accountid = req.body.accountid as string;
  const amount = req.body.amount as number;

  if (!authService.allowed(req, "accounts_credit")) {
    res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
    return;
  }
  const a = bankService.findAccountByID(accountid);
  if (a === null) {
    res.status(500).json({ errorcode: 500, errormessage: "No such account" });
    return;
  }
  if (req.userRole === Role.User && req.userID !== a.userid) {
    res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
    return;
  }

  const balance = bankService.credit(a, amount);

  res.status(200).json({ accountid: a.id, balance });
};

export const debit = (req: Request, res: Response) => {
  const accountid = req.body.accountid as string;
  const amount = req.body.amount as number;

  if (!authService.allowed(req, "accounts_debit")) {
    res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
    return;
  }
  const a = bankService.findAccountByID(accountid);
  if (a === null) {
    res.status(500).json({ errorcode: 500, errormessage: "No such account" });
    return;
  }
  if (req.userRole === Role.User && req.userID !== a.userid) {
    res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
    return;
  }
  if (a.balance < amount) {
    res.status(500).json({ errorcode: 500, errormessage: "Not enough funds" });
    return;
  }
  const balance = bankService.debit(a, amount);

  res.status(200).json({ accountid: a.id, balance });
};
