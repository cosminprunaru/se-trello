<template>
 <div class="list">    
  <!-- Hack: this will throw a null error at first before rendering -->
  <h1 v-if="currentBoard != null">{{ currentBoard.name }}</h1>
  <router-link class="btn btn-success" to="/">Home</router-link>
  <div>
    <b-btn v-b-modal.modalCreateList>Create List</b-btn>

    <!-- Modal Component -->
    <b-modal id="modalCreateList" 
             title="Create new list"
             ref="newListModal"
             @ok="handleOkCreateList"
             @shown="clearNewListName" >

      <form @submit.stop.prevent="handleSubmit">
        <b-form-input type="text"
                      ref="focusNewListInput"
                      placeholder="Enter new list name"
                      v-model="newListName"></b-form-input>
       </form>
    </b-modal>
  </div>
  <br>
    <section>
      <div
        v-for="(list, index) in all_lists"
        :value="list"
        :key="index"
        :id="list.id"
      >
          <div class="column">
            <b-card>
                     
              <span v-if="!list.isEditing">{{ list.name }}</span>
              <input v-else type="text" v-model="editingListName" placeholder="New list name here..." class="form-control">

              <br>
              <!-- Insert cards here -->
              <table class="table text-left">
                <thead>
                  <tr>
                    <!-- <th scope="col">#</th> -->
                    <th scope="col">Name</th>
                    <th scope="col">Description</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="(card, cardIndex) in list.cards"
                    :value="card"
                    :key="cardIndex"
                    :id="card.id"
                  >
                    <td> 
                      <span v-if="!card.isEditing">{{ card.name }}</span>
                      <input v-else type="text" v-model="editingCardName" placeholder="..." class="form-control">
                    </td>
                    <td> 
                      <p v-if="!card.isEditing">{{ card.description }}</p>
                      <input v-else type="text" v-model="editingCardDescription" placeholder="..." class="form-control">
                    </td>
                    <td> 
                      <button v-if="!card.isEditing" title="Edit list" class="btn btn-warning" v-on:click="setCardEditable(list, card)">Edit</button>
                      <button v-if="card.isEditing" title="Save list" class="btn btn-success" v-on:click="editCard(list, card)">Save</button>
                    </td>
                    <td> 
                      <button title="DeleteCard" class="btn btn-danger" v-on:click="deleteCard(list, card)">X</button>
                    </td>
                  </tr>
                </tbody>
              </table>
              <div class="row">
                <table class="table text-left">
                  <tbody>
                    <tr>
                      <td> 
                        <input type="text" v-model="newCardName[list.id]" placeholder="Name..." class="form-control">
                      </td>
                      <td> 
                        <input type="text" v-model="newCardDescription[list.id]" placeholder="Description..." class="form-control">
                      </td>
                      <td class="text-right">
                        <button title="New card" class="btn btn-success" v-on:click="addCard(list)">Add card</button>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
              <br>
              <button v-if="!list.isEditing" title="Edit list" class="btn btn-warning" v-on:click="setEditable(list)">Edit list</button>
              <button v-if="list.isEditing" title="Save list" class="btn btn-success" v-on:click="editList(list)">Save list</button>
              <div class="divider"/>
              <b-button v-on:click="deleteList(list)" class="btn btn-danger">Delete list</b-button>
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
  name:'List',
  data() {
    return {
      newListName: '',
      editingListName: '',
      editingCardName: '',
      editingCardDescription: '',
      currentBoard: null,
      newCardName: [],
      newCardDescription: [],
      api: 'http://localhost:3000/boards/',
      all_lists: []
    }
  },
  mounted() {
    axios
      .get(this.api + this.$route.params.id, { crossdomain: true })
      .then(response => {
        this.currentBoard = response.data
        this.all_lists = response.data.lists
      })
      .catch(error => {
        // eslint-disable-next-line
        console.log(error)
      })
      this.initEdit()
      this.initEditCard()
  },
  methods: {
    initEdit() {
      this.editingListName = ''
      for (var index = 0; index < this.all_lists.length; ++index) {
        this.all_lists[index].isEditing = false;
      }
    },
    initEditCard() {
      this.editingCardName = ''
      this.editingCardDescription = ''
      for (var index = 0; index < this.all_lists.length; ++index) {
        for (var index2 = 0; index2 < this.all_lists[index].cards.length; ++index2) {
          this.all_lists[index].cards[index2].isEditing = false;
        }
      }
    },
    setEditable(list) {
      this.editingListName = list.name
      for (var index = 0; index < this.all_lists.length; ++index) {
        if (this.all_lists[index].id === list.id) {
          this.all_lists[index].isEditing = true
          this.$forceUpdate()
          return
        }
      }
    },
    setCardEditable(list, card) {
      this.editingCardName = card.name
      this.editingCardDescription = card.description
      for (var index = 0; index < this.all_lists.length; ++index) {
        if (list.id === this.all_lists[index].id) {
          for (var index2 = 0; index2 < this.all_lists[index].cards.length; ++index2) {
            if (this.all_lists[index].cards[index2].id === card.id) {
              this.all_lists[index].cards[index2].isEditing = true;
              this.$forceUpdate()
              return
            }
          }
        }
      }
    },
    addList(lName) {
      axios
        .post(this.api + this.$route.params.id + '/lists',{
          name: lName,
          cards: []
        })
        .then(response => {
          this.currentBoard = response.data
          this.all_lists = response.data.lists
        })
        .catch(error => {
          // eslint-disable-next-line
          console.log(error)
        })
      this.initEdit()
      this.initEditCard()
    },
    editList(list) {
      axios
        .put(this.api + this.$route.params.id + '/lists',{
          id: list.id,
          name: this.editingListName,
          cards: list.cards
        })
        .then(response => {
          this.currentBoard = response.data
          this.all_lists = response.data.lists
        })
        .catch(error => {
          // eslint-disable-next-line
          console.log(error)
        })
      this.initEdit()
      this.initEditCard()
    },
    deleteList(list) {
      axios
        .delete(this.api + this.$route.params.id + '/lists/' + list.id)
        .then(response => {
          this.currentBoard = response.data
          this.all_lists = response.data.lists
        })
        .catch(error => {
          // eslint-disable-next-line
          console.log(error)
        })
      this.initEdit()
      this.initEditCard()
    },
    addCard(list) {
      axios
        .post(this.api + this.$route.params.id + '/lists/' + list.id + '/cards',{
          name: this.newCardName[list.id],
          description: this.newCardDescription[list.id]
        })
        .then(response => {
          this.currentBoard = response.data
          this.all_lists = response.data.lists
        })
        .catch(error => {
          // eslint-disable-next-line
          console.log(error)
        })
      this.newCardName[list.id] = ''
      this.newCardDescription[list.id] = ''
      this.initEdit()
      this.initEditCard()
    },
    editCard(list, card) {
      axios
        .put(this.api + this.$route.params.id + '/lists/' + list.id + '/cards',{
          id: card.id,
          name: this.editingCardName,
          description: this.editingCardDescription
        })
        .then(response => {
          this.currentBoard = response.data
          this.all_lists = response.data.lists
        })
        .catch(error => {
          // eslint-disable-next-line
          console.log(error)
        })
        this.initEdit()
        this.initEditCard()
    },
    deleteCard(list, card) {
      axios
        .delete(this.api + this.$route.params.id + '/lists/' + list.id + '/cards/' + card.id)
        .then(response => {
          this.currentBoard = response.data
          this.all_lists = response.data.lists
        })
        .catch(error => {
          // eslint-disable-next-line
          console.log(error)
        })
      this.initEdit()
      this.initEditCard()
    },
    clearNewListName () {
      this.newListName = ''
      this.$refs.focusNewListInput.focus()
    },
    handleOkCreateList (evt) {
      // Prevent modal from closing
      evt.preventDefault()
      if (!this.newListName) {
        alert('List name cannot be empty!')
      } else {
        this.handleSubmit()
      }
    },
    handleSubmit () {
      this.addList(this.newListName)
      this.clearNewListName()
      this.$refs.newListModal.hide()
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
