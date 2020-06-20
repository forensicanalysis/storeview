<template>
  <v-container class="task" fluid>
    <v-row class="pt-3" v-for="(item, i) in tasks" :key="i">
      <v-col cols="2" class="py-0">
        <b>{{ item.text }}</b>
      </v-col>
      <v-col class="py-0">
        <task class="ml-3 flex-grow-1" :name="item.text"/>
      </v-col>
    </v-row>
    <v-row class="justify-center pa-5">
      <v-btn x-large color="primary">
        <v-icon>mdi-play</v-icon>
        Run
      </v-btn>
    </v-row>
  </v-container>
</template>

<script>
  import {invoke} from "../store/invoke";
  import task from '../components/Task.vue';

  export default {
    name: 'tasks',
    components: {
      task,
    },
    data() {
      return {
        tasks: [],
      };
    },

    methods: {
      loadTasks() {
        invoke('GET', '/tasks', [], (data) => {
          this.tasks = [];
          for (let i = 0; i < data.length; i += 1) {
            this.tasks.push({
              text: data[i],
            })
          }
        });
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

  .task .v-input--checkbox {
    margin-top: 0;
  }
</style>
