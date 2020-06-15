/* tslint:disable:max-classes-per-file */

export enum Role {
  BankManager = "BankManager",
  User = "User",
}

export interface AccountListDisplay {
  name: string;
  accountid: string;
  balance: number;
  createdat: Date;
}

export interface UserList {
  id: string;
  fullname: string;
  username: string;
  createdat: Date;
  numaccounts: number;
}

export interface Transaction {
  amount: number;
  operation: string;
  balance: number;
  createdat: Date;
}
export interface Account {
  id: string;
  userid: string;
  createdat: Date;
  balance: number;
  activity: Transaction[];
}

export interface User {
  userid: string;
  name: string;
  username: string;
  password: string;
  role: Role;
  createdat: Date;
  accounts: Account[];
}

interface StrStrMap {
  [index: string]: string;
}

interface StrIntMap {
  [index: string]: number;
}
export interface Bank {
  version: number;
  accountserial: number;
  accountscatalog: StrStrMap;
  nextuserid: number;
  useridcatalog: StrIntMap;
  usernamecatalog: StrIntMap;
  users: User[];
}

interface UserService {
  create: (a: string, b: string, c: string, d?: Role) => string;
  userList: () => UserList[];
}

export let userService: UserService;

interface BankService {
  bank: Bank;
  init: () => void;
  nextuserid: () => string;
  nextaccountid: () => string;
  flush: () => void;
  addUser: (u: User) => void;
  addAccount: (u: string, a: Account) => void;
  findUserByUserName: (a: string) => User;
  findUserByID: (a: string) => User;
  accountList: () => AccountListDisplay[];
  createAccount: (a: string, b: number) => Account;
  findAccountByID: (a: string) => Account;
  credit: (a: Account, b: number) => number;
  debit: (a: Account, b: number) => number;
}
import fs from "fs";

export let bankService: BankService;
const Version = 1;

class BSClass implements BankService {
  public bank: Bank;
  private bankfile = "data/bank.json";
  init = () => {
    if (fs.existsSync(this.bankfile)) {
      this.bank = JSON.parse(fs.readFileSync(this.bankfile, "utf8"));
    }
    if (!this.bank || this.bank.version !== Version) {
      this.bank = {
        version: Version,
        accountserial: 10001,
        accountscatalog: {} as StrStrMap,
        nextuserid: 1000000,
        useridcatalog: {} as StrIntMap,
        usernamecatalog: {} as StrIntMap,
        users: [] as User[],
      };
      // console.log(this.bank);
      const userid = userService.create(
        "Stockson Bonds",
        "bankmanager",
        "headhoncho",
        Role.BankManager
      );
      this.flush();
    }
    // console.log(this.bank);
  };

  nextuserid = (): string => {
    this.bank.nextuserid++;
    this.flush();
    return this.bank.nextuserid.toString();
  };

  nextaccountid = (): string => {
    this.bank.accountserial++;
    this.flush();
    return this.bank.accountserial.toString();
  };

  addUser = (u: User) => {
    this.bank.users.push(u);
    this.bank.useridcatalog[u.userid] = this.bank.users.length - 1;
    this.bank.usernamecatalog[u.username] = this.bank.users.length - 1;
    this.flush();
  };

  flush = () => {
    fs.writeFileSync(this.bankfile, JSON.stringify(this.bank));
  };

  findUserByID = (id: string): User | null => {
    if (!this.bank.useridcatalog[id]) {
      return null;
    }

    const ix = this.bank.useridcatalog[id];
    const u = this.bank.users[ix];
    return u;
  };

  findUserByUserName = (uname: string): User | null => {
    const ix = this.bank.usernamecatalog[uname];
    if (ix === undefined) {
      return null;
    }
    return this.bank.users[ix];
  };

  findAccountByID = (id: string): Account | null => {
    if (!this.bank.accountscatalog[id]) {
      return null;
    }
    const uid = this.bank.accountscatalog[id];
    const u = this.findUserByID(uid);
    const a = u.accounts.find((e) => e.id === id);
    return a;
  };

  addAccount = (userid: string, account: Account) => {
    const ix = this.bank.useridcatalog[userid];
    this.bank.users[ix].accounts.push(account);
    this.bank.accountscatalog[account.id] = userid;
  };

  accountList = (): AccountListDisplay[] => {
    const list: AccountListDisplay[] = [];
    this.bank.users.forEach((u) => {
      u.accounts.forEach((a) => {
        list.push({
          name: u.name,
          accountid: a.id,
          balance: a.balance,
          createdat: a.createdat,
        });
      });
    });
    return list;
  };

  createAccount = (userid: string, balance: number): Account => {
    const id = this.nextaccountid();
    const t: Transaction = {
      amount: balance,
      balance,
      createdat: new Date(),
      operation: "OB",
    };
    const a: Account = {
      id,
      userid,
      balance,
      createdat: new Date(),
      activity: [t],
    };
    this.bank.accountscatalog[id] = userid;
    const u = this.findUserByID(userid); // user exists, we checked in the controller
    u.accounts.push(a);
    this.flush();
    return a;
  };

  credit = (account: Account, amount: number): number => {
    const balance = account.balance + amount;
    const t: Transaction = {
      amount,
      balance,
      createdat: new Date(),
      operation: "CR",
    };
    account.balance += amount;
    account.activity.push(t);
    this.flush();

    return t.balance;
  };

  debit = (account: Account, amount: number): number => {
    const balance = account.balance - amount;
    const t: Transaction = {
      amount,
      balance,
      createdat: new Date(),
      operation: "DB",
    };
    account.balance -= amount;
    account.activity.push(t);
    this.flush();
    return t.balance;
  };
}

class USClass implements UserService {
  create = (
    fullname: string,
    username: string,
    password: string,
    role = Role.User
  ): string => {
    // console.log("kkkkkk", bankService.bank);

    const userid = bankService.nextuserid();
    const u: User = {
      userid: userid.toString(),
      name: fullname,
      username,
      password,
      role,
      createdat: new Date(),
      accounts: [] as Account[],
    };
    bankService.addUser(u);
    return u.userid;
  };

  userList = (): UserList[] => {
    const list: UserList[] = bankService.bank.users
      .map((u: User) => {
        if (u.role !== Role.BankManager) {
          return {
            id: u.userid,
            fullname: u.name,
            username: u.username,
            createdat: u.createdat,
            numaccounts: u.accounts.length,
          };
        }
      })
      .filter((l) => l != null);
    return list;
  };
}

userService = new USClass();
bankService = new BSClass();
bankService.init();
