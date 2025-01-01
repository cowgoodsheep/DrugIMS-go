var express = require('express');
var {register,login,repairUserInfo}= require('../api/index');
var router = express.Router();
router.post('/register', async(req, res, next) => {
    const {user_name,password,role} = req.body;
    const data =await register(user_name,password,role)
    console.log(data)
    res.json(data);
  });
  router.post('/login', async(req, res, next) => {
    const {user_name,password} = req.body;
    const data =await login(user_name,password)
    res.json(data);
  });


  module.exports = router;