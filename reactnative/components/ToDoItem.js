import React, { Component } from 'react';
import { AppRegistry, FlatList, StyleSheet, Text, View, TouchableOpacity } from 'react-native';


export default class ToDoItem extends Component {
   text = ""

   render(){
       return(
           <View key={this.props.keyval} style = {styles.item}>
               <Text style={styles.itemText}>{this.props.val.item}</Text>
               <TouchableOpacity onPress={this.props.deleteMethod} style={styles.itemDelete}>
                   <Text style={styles.itemDeleteText}>X</Text>
               </TouchableOpacity>
           </View>
       );
   }

}

const styles = StyleSheet.create({
   item: {
       position: 'relative',
       padding: 20,
       paddingRight: 100,
   },
   itemText:{
       paddingLeft: 20,
       borderLeftWidth: 5,
       borderLeftColor: '#F5FCFF',
       fontSize:18,
   },
   itemDelete: {
       position:'absolute',
       justifyContent:'center',
       alignItems: 'center',
       backgroundColor: '#abacad',
       padding: 10,
       top: 10,
       bottom: 10,
       right: 10,
   },
   itemDeleteText:{
       color: 'white',
   }
});
