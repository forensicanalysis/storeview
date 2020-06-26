<template>
  <v-container class="task" fluid>
    <v-row class="pt-3" v-for="(task, name) in tasks" :key="name">
      <v-col cols="2" class="py-0">
        <v-checkbox
          v-model="active[name]"
          :label="name"
        />
      </v-col>
      <v-col class="py-0">
        <v-form ref="form" v-model="valids[name]">
          <v-jsf v-model="models[name]" :schema="cleanSchema(task)" :options="options"/>
        </v-form>
      </v-col>
    </v-row>
    <v-row class="justify-center pa-5">
      <v-btn x-large color="primary" @click="run">
        <v-icon>mdi-play</v-icon>
        Run
      </v-btn>
    </v-row>
  </v-container>
</template>

<script>
  import {invoke} from "../store/invoke";
  import VJsf from '@koumoul/vjsf';

  export default {
    name: 'tasks',
    components: {
      VJsf
    },
    data() {
      return {
        tasks: [],
        active: {},
        valids: {},
        models: {},
        options: {
          debug: false,
          disableAll: false,
          autoFoldObjects: true
        },
      }
    },

    methods: {
      loadTasks() {
        invoke('GET', '/tasks', [], (tasks) => {
          this.tasks = tasks;
        });
      },

      cleanSchema(schema) {
        delete schema["properties"]["add-to-store"];
        delete schema["properties"]["format"];
        delete schema["properties"]["output"];
        delete schema["properties"]["filter"]; // TODO
        return schema
      },

      run() {
        let that = this;
        that._.forEach(this.tasks, function (task, name) {
            if (that.active[name]) {
              invoke('POST', '/run?name=' + name, that.models[name], (result) => {
                console.log(result);
              });
            }
          }
        );
      },
    },

    mounted() {
      this.loadTasks();
    },
  };
</script>

<style>
  .task > .row {
    border-top: 1px solid rgba(0, 0, 0, 0.12);
  }

  .task .vjsf-property .row {
    flex-direction: row;
    flex-wrap: nowrap;
  }

  .task .col-2 label {
    font-weight: bold;
    color: #333;
    font-size: 1rem;
  }

  .task .v-input--checkbox {
    margin-top: 0;
  }
</style>
