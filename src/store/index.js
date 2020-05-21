import Vue from 'vue';
import Vuex from 'vuex';
import { invoke } from './invoke';
import { templates } from './templates';

Vue.use(Vuex);

export default new Vuex.Store({

  state: {
    tables: [],

    type: 'file',
    items: [],
    all: [],
    header: [],
    sort: {
      table: '',
      columns: {}
    },
    filter: {
      table: '',
      columns: {}
    },
    limit: 30,
    offset: 0,
    templates: templates,

    item: {},
    listPane: 80,

    tasks: {
      'fetch': {},
      'unpack': { 'requires': ['fetch'] },
      'plaso': { 'requires': ['unpack'] },
      'hash': { 'requires': ['unpack'] },
      'compare': { 'requires': ['hash'] },
      'alert': { 'requires': ['plaso', 'compare'] },
    },

  },

  mutations: {

    setLimit(state, data) {
      state.limit = data;
    },

    setOffset(state, data) {
      state.offset = data;
    },

    setTable(state, data) {
      state.sort = {
        table: '',
        columns: {}
      };
      state.filter = {
        table: '',
        columns: {}
      };
      state.type = data;
    },

    setTables(state, data) {
      state.tables = data;
    },

    setItems(state, data) {
      let calcHeaders = function () {
        let headers = [];
        let keys = _.keys(data[0]);
        for (let i = 0; i < keys.length; i++) {
          headers.push({
            text: Vue._.startCase(keys[i]),
            value: keys[i]
          });
        }
        return headers;
      };

      if (data.length === 0) {

      } else if ('type' in data[0]) {
        state.type = data[0]['type'];
        if (data[0]['type'] in state.templates && 'headers' in state.templates[data[0]['type']]) {
          state.headers = state.templates[data[0]['type']]['headers'];
        } else {
          state.headers = calcHeaders();
        }
      } else {
        state.headers = calcHeaders();
      }
      state.offset = 0;
      state.items = data;
      console.log('SIZE');
      console.log(state.items.length);
    },

    setItem(state, data) {
      if (Vue._.isEmpty(data)) {
        state.listPane = 80;
      } else {
        if (state.listPane >= 99.9) {
          state.listPane = 50;
        }
      }
      state.item = data;
    },

    setlistPane(state, data) {
      state.listPane = data;
    },
    setSort(state, data) {
      state.sort = data;
    },

    setFilter(state, data) {
      state.filter = data;
    },

  },

  actions: {

    created({ commit, dispatch, state }) {
      invoke('GET', '/tables', [], function (tables) {
        commit('setTables', tables);
      });

      dispatch('loadItems');

    },

    loadItems({ commit, state }) {
      let url = '/items?type=' + state.type;
      console.log(state.type)
      if (!Vue._.isEmpty(state.filter) && state.filter.type !== '' && !Vue._.isEmpty(state.filter.columns)) {
        Vue._.forEach(state.filter.columns, function (value, column) {

          for (let i=0; i<value.length; i++) {
            url += '&filter[' + column + ']=' + encodeURI(value[i]);
          }

          console.log(url)
        });
      }
      if (!Vue._.isEmpty(state.sort) && state.sort.type !== '' && !Vue._.isEmpty(state.sort.columns)) {
        Vue._.forEach(state.sort.columns, function (value, column) {
          let direction = '';
          switch (value) {
            case 'ASC':
              direction = 'ASC';
              break;
            case 'DESC':
              direction = 'DESC';
          }
          if (direction !== '') {
            url += '&sort[' + column + ']=' + direction;
          }
        });
      }

      url += '&offset=' + state.offset;
      url += '&limit=' + state.limit;

      console.log(url);

      invoke('GET', url, [], function (items) {
        console.log(items)
        commit('setItems', items);
      });

    },

    loadDirectories({ commit, state }, payload) {
      return new Promise(resolve => {

        let directories = []
        let slash = ''
        let url = ''

        if (state.type === 'directory') {
          slash = '/'
        }
        else if (state.type === 'file') {
          slash = '/'
        }
        else if (state.type === 'windows-registry-key') {
          slash = '\\'
        }
        else {
          console.log("TABLE DOES NOT EXIST!")
        }

        if((state.type === 'windows-registry-key') && (payload.path === '')) {
          url = '/tree?directory=' + payload.path;
        }
        else {
          url = '/tree?directory=' + payload.path + slash;
        }

        url += '&type=' + state.type

        // url = '/tree?directory=/C/&type=directory'
        console.log('***************\n', url)

        invoke('GET', url, [], function (data) {
          console.log(data)
          return
        }).then(response => {

          for (let i = 0; i < response.length; i++) {

            if((state.type === 'windows-registry-key') && (payload.path === '')) {
              directories.push(
                {
                  path: response[i],
                  name: response[i],
                  children: []
                }
              );
            }
            else {
              directories.push(
                {
                  path: payload.path + slash + response[i],
                  name: response[i],
                  children: []
                }
              );
            }
            resolve(directories)
          }

        })

      });
    },

    updateSort({ commit, dispatch }, data) {
      commit('setSort', data);
      dispatch('loadItems');
    },

    updatePage({ commit, dispatch }, data) {
      commit('setOffset', data);
      dispatch('loadItems');
    },
  },

  modules: {},

});
