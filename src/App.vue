<template>
  <v-app>
    <v-app-bar app color="primary" dark clipped-right clipped-left flat>
      <v-toolbar-title class="mr-8">Elementary</v-toolbar-title>
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
  body {
    font-family: 'Roboto', sans-serif;
  }

  h1, h2, h3, .v-toolbar {
    font-family: 'Roboto Slab', serif;
  }

  .v-btn--fab {
    border-radius: 0;
    border-radius: 0;
  }

  dt {
    font-weight: bold;
  }

  dd {
    margin-left: 12px;
  }
</style>
