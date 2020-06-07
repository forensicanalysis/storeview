<template>
  <v-list dense>
    <v-subheader>Tasks</v-subheader>
    <v-list-item-group v-model="tasks" color="primary">
      <v-list-item v-for="(item, i) in tasks" :key="i">
        <v-list-item-icon>
          <v-icon>mdi-play</v-icon>
        </v-list-item-icon>
        <v-list-item-content>
          <v-list-item-title v-text="item.text"></v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list-item-group>
  </v-list>
</template>

<script>
  import {invoke} from "../store/invoke";

  export default {
    name: 'tasks',
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
              // icon:
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
