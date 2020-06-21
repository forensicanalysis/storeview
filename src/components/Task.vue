<template>
  <div>
    <v-form ref="form" v-model="valid">
      <v-jsf v-model="model" :schema="schema" :options="options"/>
    </v-form>
    <!--v-layout row class="mt-2">
      <v-spacer/>
      <v-btn color="primary" @click="$refs.form.validate()">Validate</v-btn>
    </v-layout-->
  </div>
</template>

<script>
  import {invoke} from "../store/invoke";
  import VJsf from '@koumoul/vjsf';

  export default {
    name: 'task',
    props: ['name'],
    components: {
      VJsf
    },
    data() {
      return {
        valid: true,
        model: {},
        schema: {},
        options: {
          debug: false,
          disableAll: false,
          autoFoldObjects: true
        },
      };
    },
    methods: {
      loadSchema() {
        invoke('GET', '/task?name=' + this.name, [], (data) => {
          this.schema = data;
        });
      },
    },

    mounted() {
      this.loadSchema();
    },
  };
</script>
