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
  <v-container class="task scrollableArea" fluid>
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
