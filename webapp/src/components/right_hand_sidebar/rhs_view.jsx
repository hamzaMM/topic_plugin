import React from 'react';
import PropTypes from 'prop-types';
import axios from 'axios';
import ListGroup from 'react-bootstrap/ListGroup'
import './list.css'

export default class RHSView extends React.PureComponent {
    static propTypes = {
        user: PropTypes.object.isRequired,
        channel: PropTypes.array.isRequired,
        posts: PropTypes.array.isRequired,
    }
    constructor(props){
        super(props);    // pass props to "father" constructor
        this.state = {
          topics : null,
          input: '',
          list: []
        }

    this.handleAdd = this.handleAdd.bind(this);
    this.dltHandler = this.deleteHandler.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
      }

      onListChange(newList) {
        this.setState({
          list: newList
        });
      }
      handleAdd() {
      if(this.state.input != ''){
      this.setState({list: [...this.state.list, this.state.input]})
      }
      this.setState({input: ''})
      }
      handleSubmit() {
        console.log("running ")
        const url = "/plugins/topic_modeling/add_topics";
        var temp = this.state.list.toString()
        temp = '"' + temp.split( "," ).join( '","' ) + '"';
        console.log("sending to submit")
        console.log(temp)
        this.setState({
          topics: temp
        });        
        axios({
            url: url,
            method: 'post',
            data: JSON.stringify({labels: temp , user: this.props.user})
                    })
          .catch(err => console.log(err));
      }

      deleteHandler(listitem){
        var index = this.state.list.indexOf(listitem)
        var temp = this.state.list.filter((data, idx) => idx !== index )
        this.setState({list: temp})
      }
    // AXIOS request

    componentDidMount() {
        const url = "/plugins/topic_modeling/topics";
        axios({
            url: url,
            method: 'post',
            data: this.props.user
            })
        .then((response) => {
            var array = response.data.replace(/["]+/g, '')
            array = array.split(',');
            this.setState({list: array})
            console.log("array after mounting:")
            console.log(this.state.list)
            })
          .catch(err => console.log(err));
        }    


    render() {
              
        return (
        <div>
      <ListGroup  className ="list-group-mine" defaultActiveKey="">
      {this.state.list.map(listitem => (
    <ListGroup.Item action onClick={this.deleteHandler.bind(this, listitem)}>
      {listitem}
    </ListGroup.Item>         
     ))}
  </ListGroup>,
      <input value={this.state.input} onInput={e =>  this.setState({input: e.target.value})}/>
      <button type="button" className="btn btn-primary" onClick={this.handleAdd}>Add</button>
      <br />
      <button type="button" className="btn btn-primary" onClick={this.handleSubmit}>Save</button>

          {console.log(this.state.topics)}
        </div>
            
        );
    }
}




