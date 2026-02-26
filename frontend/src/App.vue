<template>
  <div id="app-container">
    <router-view />
  </div>
</template>

<style>
/* Basic reset since we aren't using Tailwind */
#app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}
</style>