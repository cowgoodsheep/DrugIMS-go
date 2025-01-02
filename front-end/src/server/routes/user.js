
var express = require('express');
var {updateUserInfo,getUserList,deleteUser,updateUser}= require('../api/index');
var router = express.Router();
router.post('/setUserInfo', async(req, res, next) => {
    const {user_name,password,telephone,address,user_id} = req.body;
    const data =await updateUserInfo(user_name,password,telephone,address,user_id)
    res.json(data);
  });
  router.post('/getUserList', async(req, res, next) => {
    const {searchValue} = req.body
    const data =await getUserList(searchValue)
    res.json(data);
  });
  router.post('/deleteUser', async(req, res, next) => {
    const {user_id} = req.body
    const data =await deleteUser(user_id)
    res.json(data);
  });
  router.post('/updateUser', async(req, res, next) => {
    const {user_name,password,role,telephone,address,user_id} = req.body
    const data =await updateUser(user_name,password,role,telephone,address,user_id)
    res.json(data);
  });
  module.exports = router;