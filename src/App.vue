<!--
Copyright (c) 2020 Siemens AG

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

Author(s): Jonas Plum
-->
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
