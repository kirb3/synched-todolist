import React, { Component } from 'react';
import {Platform,StyleSheet,Text,View,Image,ImageBackground,TextInput,ScrollView,TouchableOpacity,AsyncStorage,RefreshControl} from 'react-native';
import ToDoItem from './components/ToDoItem'

export default class App extends Component {

    constructor(props) {
        super(props);
        this.state = {
            refreshing: false,
            itemArray: [{'item':'item1'}],
            itemText: '',
        };
    }
    _onRefresh() {
        this.setState({refreshing: true});
        this.getListCloud().then(() => {
            super.setState({refreshing:false});
        })
    }

    componentDidMount(){
        this.getListCloud();
    }

    async saveListCloud(){ //Update the list in server
        var itemArray = this.state.itemArray;
        var sendList = [];
        itemArray.forEach(function(items){
            sendList.push(items.item);
        });

        try {            
            let response = await fetch('http://localhost:3000/user/hello2',
            {
                method: 'PUT',
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    "username":"hello2",
                    "items":sendList
                }),
            });
            let responseMessage = await response;
            //console.log('response: '+JSON.stringify(responseMessage, null, 4));
        } catch (error) {
            console.log(error);
        }
    }
    async getListCloud(){ //Get the list in the server
        try {
            let response = await fetch('http://localhost:3000/user/hello2');
            let responseJson = await response.json();
            console.log(responseJson.items);
            var itemList = responseJson.items.map(function(value) {
                return {'item':value}
            });
            this.setState({itemArray:itemList});
        } catch (error) {
            console.error(error);
        }
    }
    async createListCloud(){ //Create new list in the server
        try {
            let response = await fetch('http://localhost:3000/user',
            {
                method: 'POST',
                headers: {
                    Accept: 'application/json',
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    "username": "hello2",
                    "items": ["cool","updated"]
                }),
            });
            // let responseJson = await response.json();
            console.log(response);
        } catch (error) {
            console.error(error);
        }
    }

    addItem(){
        if (this.state.itemText) {
            this.state.itemArray.push({'item':this.state.itemText});
            this.setState({itemArray:this.state.itemArray});
            //this.saveListLocal(this.state.itemArray);
            this.setState({itemText:''});
            this.saveListCloud();
        }
    }

    deleteItem(key){
        this.state.itemArray.splice(key, 1);
        this.setState({itemArray:this.state.itemArray});
        //this.saveListLocal(this.state.itemArray);
        this.saveListCloud();
    }

    render() {

        let items = this.state.itemArray.map((val,key) => {
            return <ToDoItem key={key} keyval={key} val={val} deleteMethod = {() => this.deleteItem(key)} />
        });

        return (
            <View style={styles.container}>

                <View style={styles.header}>
                    <Text style={styles.headerText}> To Do List</Text>
                </View>

                <ScrollView 
                    refreshControl={
                        <RefreshControl
                            refreshing={this.state.refreshing}
                            onRefresh={this._onRefresh.bind(this)}
                        />
                    }
                    style={styles.scrollContainer}>
                {items}
                </ScrollView>

                <View style={styles.footer}>
                    <TouchableOpacity onPress={this.addItem.bind(this)} style={styles.addItemButton}>
                        <Text style={styles.addItemButtonText}>+</Text>
                    </TouchableOpacity>

                <TextInput
                    style={styles.textInput}
                    onChangeText={(itemText) => this.setState({itemText})}
                    value={this.state.itemText}
                    placeholder='[add task here]'
                    placeholderTextColor ='white'
                    underlineColorAndroid='transparent'>
                </TextInput>

                </View>

            </View>
        );
    }
}

const styles = StyleSheet.create({

 container: {
   flex: 1,
   justifyContent: 'center',
   alignItems: 'center',
   backgroundColor: '#F5FCFF',
 },
 header:{
   backgroundColor:'#1578a3',
   alignItems:'center',
   borderTopWidth:15,
   borderTopColor:'#F5FCFF',
   borderBottomWidth:10,
   borderBottomColor:'#F5FCFF',
 },
 headerText:{
   color:'white',
   fontSize: 18,
   padding: 26,
 },
 scrollContainer:{
   flex:1,
   marginBottom: 100,
 },
 footer: {
   position: 'absolute',
   alignItems: 'center',
   bottom: 0,
   left: 0,
   right: 0,
 },
 addItemButton:{
   backgroundColor: '#1578a3',
   width:90,
   height: 90,
   borderRadius: 50,
   borderColor: '#F5FCFF',
   alignItems: 'center',
   justifyContent: 'center',
   elevation: 10,
   marginBottom: -45,
   zIndex: 10, 
 },
 addItemButtonText:{
   color: '#fff',
   fontSize: 24,
 },
 textInput: {
   alignSelf: 'stretch',
   color: '#fff',
   padding: 20,
   paddingTop: 46,
   backgroundColor: '#252525',
   borderTopWidth: 2,
   borderTopColor: '#ededed',
 },
 img:{
   flex:1,
 }
});
