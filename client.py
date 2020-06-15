#!/usr/local/bin/python3

import getpass
import json
import os
import pprint
import random
import re
import sys
import urllib
from datetime import datetime

import httplib2

pp = pprint.PrettyPrinter(indent=3)
sessions = {
    'bankmanager' : {'username': 'bankmanager', 'password': 'headhoncho', 'token': ''},
    'a1003' : {'username': 'a1003', 'password': 'ppp', 'token': ''},
    'user4' : {'username': 'user4', 'password': 'pass', 'token': ''}
}

currentuser = ''
headers = {'Content-type': 'application/json'}
debug = False

    
def readinput(prompt='Enter Input'):
    sys.stdout.write(prompt+": ")
    sys.stdout.flush()
    inp = sys.stdin.readline(80)
    return inp.strip()

def sendrequest(method, endpoint, data=None):
    print(method, endpoint, data)
    if (debug):
      httplib2.debuglevel = 32
    headers['Authorization'] = sessions[currentuser]['Authorization']
    body = None
    if (data is not None):
      body = json.dumps(data)
    response, content = http.request(HOST+endpoint, method, headers=headers, body=body)
    if (debug):
      print(response, content)
    if response['status'] != '200':
        # c = json.loads(content)
        print('Error: ', content)
        return {}
    t = {}
    if (content != b''):
        t = json.loads(content)
    return t

def login(username, password):
  global currentuser
  data = {'username':username, 'password': password}
  response, content = http.request(HOST+'/login', 'POST', headers=headers, body=json.dumps(data))
  if (debug):
    print(response, content)
  if response['status'] != '200':
      # c = json.loads(content)
      print('Error: ', content)
      return {}
  t = {}
  if (content != b''):
      t = json.loads(content)
  currentuser = username
  sessions[username]['token'] = t['token']
  sessions[username]['Authorization'] = 'Bearer {0}'.format(t['token'])
  return t

def getaccounts():
  return sendrequest('GET', '/accounts')

def getusers():
  return sendrequest('GET', '/users')

def getuser(id):
  return sendrequest('GET', '/users/{0}'.format(id))

def getaccount(id):
  return sendrequest('GET', '/accounts/{0}'.format(id))

def createuser(fullname, username, password):
  return sendrequest('POST', '/users', {'username':username, 'password': password, 'fullname': fullname})

def createaccount(userid, balance):
  return sendrequest('POST', '/accounts', {'userid':userid, 'balance': balance})

def credit(accountid, amount):
  return sendrequest('POST', '/accounts/credit', {'accountid':accountid, 'amount': amount})
  
def debit(accountid, amount):
  return sendrequest('POST', '/accounts/debit', {'accountid':accountid, 'amount': amount})
  
def printobj(obj):
      print(json.dumps(obj, indent=2))

def quit():
    print("Quitting...")
    sys.exit(0)

def runtests(host):
    global debug, HOST
    HOST = host
    # username, password = 'bankmanager', 'headhoncho'
    username, password = 'user4', 'pass'
    login(username, password)
    #u = login()
    accounts = getaccounts()
    # printobj(accounts)
    u = getusers()
    user7 = getuser('1000007')
    printobj(user7)
    user8 = getuser('1000008')
    printobj(user8)
    # # printobj(u)
    # users = []
    # if 'users' in u:
    #   users = u['users']
    # printobj(users)
    # # if (len(users) > 0):
    #   rv = random.randint(0, len(users)-1)
    #   user = users[rv]
    # # # debug = True
    # # userid = createuser('User Name2', 'user4', 'pass')
    # user = getuser('1000008')
    # user = getuser('1000002')
    # printobj(user)
    # # debug = False
    v = createaccount('1000007', 500000)
    if 'accountid' in v:
      actid7 = v['accountid']
    else:
      actid7 = '10014'

    v = createaccount('1000008', 600000)
    if 'accountid' in v:
      actid8 = v['accountid']
    else:
      actid8 = '10015'
    
    credit(actid7, 500000)
    print('balance is 600000')
    printobj(getaccount(actid8))
    credit(actid8, 500000)
    print('balance is 1100000')
    printobj(getaccount(actid8))

    debit(actid7, 500000)
    debit(actid7, 500000)
    debit(actid7, 500000)

    print('balance is 1100000')
    printobj(getaccount(actid8))
    debit(actid8, 500000)
    print('balance is 600000')
    printobj(getaccount(actid8))
    debit(actid8, 500000)
    print('balance is 100000')
    printobj(getaccount(actid8))
    debit(actid8, 500000)
    print('balance is ERROR ')
    printobj(getaccount(actid8))

    # ac = user['user']['accounts']
    # rv = random.randint(0, len(ac)-1)
    # b = credit(ac[rv]['id'], 300000)
    # print (ac[rv]['id'], b['balance'])
    # rv = random.randint(0, len(ac)-1)
    # b = debit(ac[rv]['id'], 300000)
    # print (ac[rv]['id'], b['balance'])
    # b = debit(ac[rv]['id'], 300000)
    # print (ac[rv]['id'], b['balance'])
    # b = debit(ac[rv]['id'], 300000)
    # print (ac[rv]['id'], b['balance'])
    # # // pick a random number
    # # account = getaccount(id)
    # # credit(account.id, 200000)
    # # debit(account.id, 900000)

def runusertests():
    userid = createuser('User Name2', 'user4', 'pass')
    user = getuser(userid['userid'])
    printobj(user)


def usage():
        print('Usage: {0} (dev|stag|prod) '.format(sys.argv[0]))
        sys.exit(1)

def printheader(mode):
    if (mode == 'prod'):
        print('X' * 66)
        print
        print('PRODUCTION MODE')
        print
        print('X' * 66)
    return

ops = [
  #   [unapprove, 'UnApprove <corp>'],
  #   [reallocate, 'Reallocate Orders <xactid, orderids, reason, invid, invtype>'],
	# [getinvestmentdata, 'Get Investment Data'],
	# [reallocatecorp, 'Reallocate Corp <from, reason, to>'],
	# [createinvoicefile, 'Create Invoice File'],
	# [createinvoice, 'Create Invoice <corpid, description, amount>'],
  #   [sendinvoice, 'Send Invoice'],
  #   [upgradeorders, 'Upgrade Orders'],
  #   [compsub, 'Comp Subscription'],
  #   [getinvtransactions, 'Get Investment Transactions'],
  #   [createregdagreement, 'Create RegD Agreement']
]

opfns = set([s[0] for s in ops])
opfns.add(quit) 

GOHOST = 'http://localhost:9765'
NODEHOST = 'http://localhost:9765'
def main():
    global HOST, http, cookie, mode, debug
    HOST = GOHOST
    http = httplib2.Http(disable_ssl_certificate_validation=True)
    # httplib2.debuglevel = 32
    http.follow_all_redirects = True
    # runtests(GOHOST)
    # debug = True
    runtests(NODEHOST)
    # debug = False

    # # passwd = getpass.unix_getpass(prompt='FM Password: ', stream=sys.stdout)
    # t = login('0admin0', passwd)
    # while True:
    #     printheader(mode)
    #     n = 1
    #     for o in ops:
    #         print("{0}. {1}".format(n, o[1]))
    #         n = n + 1
    #     print()
    #     print('Q to Quit')
    #     o = readinput('Select an operation: ')
    #     if o == 'q':
    #         quit()
    #     op = ''
    #     o1 = int(o)-1
    #     try:
    #         op = ops[o1][0]
    #     except:
    #         print
    #         print("=================")
    #         print
    #         print("ERROR: {0} not a valid operation".format(op))
    #         print
    #         print("=================")
    #         print
    #         continue
    #     op()

if __name__ == '__main__':
    main()
