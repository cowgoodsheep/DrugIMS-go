<template>
  <div id="app">
    <nav>
      <router-link to="/home">主页</router-link>
    </nav>
    <router-view/>
  </div>
</template>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #244874;
  margin-top: 60px;
}
</style>