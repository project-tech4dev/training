"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.get = exports.create = exports.list = void 0;
const bank_1 = require("../models/bank");
const auth_1 = require("./auth");
exports.list = (req, res) => {
    if (!auth_1.authService.allowed(req, "users_list")) {
        res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
        return;
    }
    const users = bank_1.userService.userList();
    res.status(200).json({ users });
};
exports.create = (req, res) => {
    if (!auth_1.authService.allowed(req, "users_create")) {
        res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
        return;
    }
    const fullname = req.body.fullname;
    const username = req.body.username;
    const password = req.body.password;
    const user = bank_1.bankService.findUserByUserName(username);
    if (user !== null) {
        res
            .status(500)
            .json({ errorcode: 500, errormessage: "User Already Exists" });
        return;
    }
    const userid = bank_1.userService.create(fullname, username, password);
    res.status(200).json({ userid });
};
exports.get = (req, res) => {
    const userid = req.params.id;
    if (!auth_1.authService.allowed(req, "users_get")) {
        res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
        return;
    }
    if (req.userRole === bank_1.Role.User && req.userID !== userid) {
        res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
        return;
    }
    const user = bank_1.bankService.findUserByID(userid);
    if (user === null) {
        res.status(500).json({ errorcode: 500, errormessage: "No such user" });
        return;
    }
    res.status(200).json({ user });
};
//# sourceMappingURL=users.js.map