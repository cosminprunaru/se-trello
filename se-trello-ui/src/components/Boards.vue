<template>
  <div class="boards">
    <h1>Boards</h1>
    <div>
      <b-btn v-b-modal.modalCreateBoard>Create board</b-btn>

      <!-- Modal Component -->
      <b-modal id="modalCreateBoard" 
               title="Create new board"
               ref="newBoardModal"
               @ok="handleOkCreateBoard"
               @shown="clearNewBoardName" >

        <form @submit.stop.prevent="handleSubmit">
          <b-form-input type="text"
                        ref="focusNewBoardInput"
                        placeholder="Enter new board name"
                        v-model="newBoardName"></b-form-input>
        </form>
      </b-modal>
    </div>
    <br>
    <section v-if="errored">
      <p>We're sorry, we're not able to retrieve this information at the moment, please try back later</p>
    </section>

    <section v-else>
      
      <div v-if="loading">Loading...</div>
      
      <div
        v-else
        v-for="(board, index) in all_boards"
        :value="board"
        :key="index"
        :id="board.id"
      >
          <div class="column">
            <b-card>
                     
              <span v-if="!board.isEditing">
                <router-link :to="getCurrentBoard(board.id)">
                 {{ board.name }}
                </router-link>
              </span>

              <input v-else type="text" v-model="editingBoardName" placeholder="New board name here..." class="form-control">
              <br>
              <button v-if="!board.isEditing" title="Edit board" class="btn btn-warning" v-on:click="setEditable(board)">Edit</button>
              <button v-if="board.isEditing" title="Save board" class="btn btn-success" v-on:click="editBoard(board)">Save</button>
              <div class="divider"/>
              <b-button v-on:click="deleteBoard(board.name)" class="btn btn-danger">Delete</b-button>
            </b-card>
          </div>
      </div>

    </section>
  </div>
</template>

<script>

import axios from 'axios'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

export default {
  name: 'Boards',
  props: {
    msg: String
  },
  data() {
    return {
      newBoardName: '',
      api: 'http://localhost:3000',
      loading: true,
      errored: false,
      editingBoardName: '',
      all_boards: []
    }
  },
  mounted() {
    axios
      .get(this.api + '/boards', { crossdomain: true })
      .then(response => {
        this.all_boards = response.data
      })
      .catch(error => {
        // eslint-disable-next-line
        console.log(error)
      })
      .finally(() => this.loading = false)
  },
  methods: {
    getCurrentBoard(id) {
      return '/list/' + id
    },
    initEdit() {
      this.editingBoardName = ''
      if (this.all_boards == null)
        return
        
      for (var index = 0; index < this.all_boards.length; ++index) {
        this.all_boards[index].isEditing = false;
      }
    },
    setEditable(board) {
      this.editingBoardName = board.name
      for (var index = 0; index < this.all_boards.length; ++index) {
        if (this.all_boards[index].id === board.id) {
          this.all_boards[index].isEditing = true
          this.$forceUpdate()
          return
        }
      }
    },
    deleteBoard(name) {
      axios
        .delete(this.api + '/boards/' + name)
        .then(response => {
          this.all_boards = response.data
        })
        .catch(error => {
          // eslint-disable-next-line
          console.log(error)
        })
      this.initEdit()
    },
    addBoard(bName) {
      axios
        .post(this.api + '/boards',{
          name: bName,
          lists: []
        })
        .then(response => {
          this.all_boards = response.data
        })
        .catch(error => {
          // eslint-disable-next-line
          console.log(error)
        })
        this.initEdit()
    },
    editBoard(board) {
      axios
        .put(this.api + '/boards',{
          id: board.id,
          name: this.editingBoardName,
          lists: board.lists
        })
        .then(response => {
          this.all_boards = response.data
        })
        .catch(error => {
          // eslint-disable-next-line
          console.log(error)
        })

        this.initEdit()
    },
    clearNewBoardName () {
      this.newBoardName = ''
      this.$refs.focusNewBoardInput.focus()
    },
    handleOkCreateBoard (evt) {
      // Prevent modal from closing
      evt.preventDefault()
      if (!this.newBoardName) {
        alert('BoardBoard name cannot be empty!')
      } else {
        this.handleSubmit()
      }
    },
    handleSubmit () {
      this.addBoard(this.newBoardName)
      this.newBoardName = ''
      this.$refs.newBoardModal.hide()
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}

* {
  box-sizing: border-box;
}

body {
  font-family: Arial, Helvetica, sans-serif;
}

.divider{
    width:5px;
    height:5px;
    height:auto;
    display:inline-block;
}

/* Float four columns side by side */
.column {
  float: left;
  width: 25%;
  padding: 0 10px;
}

/* Remove extra left and right margins, due to padding in columns */
.row {margin: 0 -5px;}

/* Clear floats after the columns */
.row:after {
  content: "";
  display: table;
  clear: both;
}

/* Style the counter cards */
.card {
  box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2); /* this adds the "card" effect */
  padding: 16px;
  text-align: center;
  background-color: #f1f1f1;
}

/* Responsive columns - one column layout (vertical) on small screens */
@media screen and (max-width: 600px) {
  .column {
    width: 100%;
    display: block;
    margin-bottom: 20px;
  }
}
</style>
