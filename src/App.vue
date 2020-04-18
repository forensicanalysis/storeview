<template>
  <v-app>
    <v-app-bar app scolor="light-blue darken-3" style="background-color: #004c8c" dark clipped-right clipped-left>
      <v-toolbar-title class="mr-8">forensicstore</v-toolbar-title>
      <!--v-spacer/-->
      <v-toolbar-items>
        <v-btn v-for="(name, route) in menu" :key="name" @click="$router.push(route)" text>
          {{name}}
        </v-btn>
      </v-toolbar-items>
    </v-app-bar>
    <v-content>
      <v-container fluid>
        <router-view/>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>

  export default {
    name: 'app',
    components: {},
    data: function () {
      return {
        drawer: {},
        menu: {
          'items': 'Items',
          // 'workflows': 'Workflows',
        }
      };
    },
    created() {
      this.$store.dispatch('created');
    }
  };
</script>

<style>
  dt {
    font-weight: bold;
  }

  dd {
    margin-left: 12px;
  }
</style>
