// Copyright (c) 2020 Siemens AG
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// Author(s): Jonas Plum

import Vue from 'vue';
import Vuex from 'vuex';
import { invoke } from './invoke';
import { templates } from './templates';

Vue.use(Vuex);

export default new Vuex.Store({

  state: {
    tables: [],

    type: '',
    items: [],
    itemCount: 0,
    all: [],
    header: [],
    sort: {
      table: '',
      columns: {},
    },
    filter: {
      table: '',
      columns: {},
    },
    limit: 30,
    offset: 0,
    templates,

    label: '',

    item: {},
    listPane: 80,

    initialLoad: true,

    tasks: {
      fetch: {},
      unpack: { requires: ['fetch'] },
      plaso: { requires: ['unpack'] },
      hash: { requires: ['unpack'] },
      compare: { requires: ['hash'] },
      alert: { requires: ['plaso', 'compare'] },
    },

    refreshDetails: true,

  },

  mutations: {

    setType(state, data) {
      state.type = data;
    },

    setLimit(state, data) {
      state.limit = data;
    },

    setOffset(state, data) {
      state.offset = data;
    },

    setTable(state, data) {
      state.sort = {
        table: '',
        columns: {},
      };
      state.filter = {
        table: '',
        columns: {},
      };
      state.type = data;
    },

    setTables(state, data) {
      state.tables = data;
    },

    setItems(state, data) {
      const calcHeaders = function () {
        const headers = [];
        const keys = Vue._.keys(data[0]);
        for (let i = 0; i < keys.length; i += 1) {
          if (keys[i] !== "id" && keys[i] !== "type") {
            headers.push({
              text: Vue._.startCase(keys[i]),
              value: keys[i],
            });
          }
        }
        return headers;
      };

      if (data.length !== 0 && 'type' in data[0]) {
        // state.type = data[0].type;
        if (data[0].type in state.templates && 'headers' in state.templates[data[0].type]) {
          state.headers = state.templates[state.type].headers;
        } else {
          state.headers = calcHeaders();
        }
      } else {
        state.headers = calcHeaders();
      }
      state.offset = 0;
      state.items = data;
    },

    setItemCount(state, data) {
      state.itemCount = data;
    },

    setItem(state, data) {
      if (Vue._.isEmpty(data)) {
        state.listPane = 80;
      } else if (state.listPane >= 99.9) {
        state.listPane = 50;
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

    setLabelFilter(state, data) {
      state.label = data;
    },

  },

  actions: {

    async loadItems({ commit, state }) {
      let url = `/items?type=${state.type}`;
      if (!Vue._.isEmpty(state.filter) && state.filter.type !== '' && !Vue._.isEmpty(state.filter.columns)) {
        Vue._.forEach(state.filter.columns, (value, column) => {
          for (let i = 0; i < value.length; i += 1) {
            url += `&filter[${column}]=${encodeURI(value[i])}`;
          }
        });
      }
      if (!Vue._.isEmpty(state.sort) && state.sort.type !== '' && !Vue._.isEmpty(state.sort.columns)) {
        Vue._.forEach(state.sort.columns, (value, column) => {
          let direction = '';
          switch (value) {
            case 'ASC':
              direction = 'ASC';
              break;
            case 'DESC':
            default:
              direction = 'DESC';
          }
          if (direction !== '') {
            url += `&sort[${column}]=${direction}`;
          }
        });
      }
      if (state.label !== '') {
        url += `&labels=${state.label}`;
      }

      url += `&offset=${state.offset}`;
      url += `&limit=${state.limit}`;

      console.log(url);

      invoke('GET', url, [], (items) => {
        commit('setType', state.type);
        commit('setItems', items.elements);
        commit('setItemCount', items.count);
      });
    },

    async loadDirectories({ commit, state }, payload) {
      return new Promise((resolve) => {

        const directories = [];

        let slash = '';
        let url = '';

        if (state.type === 'directory') {
          slash = '/';
        } else if (state.type === 'file') {
          slash = '/';
        } else if (state.type === 'windows-registry-key') {
          slash = '\\';
        } else {
          resolve(directories);
          return
        }

        if ((state.type === 'windows-registry-key') && (payload.path === '')) {
          url = `/tree?directory=${payload.path}`;
        } else {
          url = `/tree?directory=${payload.path}${slash}`;
        }

        url += `&type=${state.type}`;

        invoke('GET', url, [], (data) => {
          return
        }).then((response) => {
          for (let i = 0; i < response.length; i += 1) {
            if ((state.type === 'windows-registry-key') && (payload.path === '')) {
              directories.push(
                {
                  path: response[i],
                  name: response[i],
                  children: [],
                },
              );
            } else {
              directories.push(
                {
                  path: payload.path + slash + response[i],
                  name: response[i],
                  children: [],
                },
              );
            }
            resolve(directories);
          }
        });
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

    updateLabel({ commit, dispatch }, data) {
      commit('setLabelFilter', data);
      dispatch('loadItems');
    },

  },

  modules: {},

});
