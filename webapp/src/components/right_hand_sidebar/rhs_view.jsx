import React from 'react';
import PropTypes from 'prop-types';
import axios from 'axios';
import {
    Accordion,
    AccordionItem,
    AccordionItemHeading,
    AccordionItemButton,
    AccordionItemPanel,
} from 'react-accessible-accordion';

// Demo styles, see 'Styles' section below for some notes on use.
import 'react-accessible-accordion/dist/fancy-example.css';

export default class RHSView extends React.PureComponent {
    static propTypes = {
        team: PropTypes.object.isRequired,
        channel: PropTypes.string.isRequired,
        posts: PropTypes.array.isRequired,
    }
    constructor(props){
        super(props);    // pass props to "father" constructor
        this.state = {
          topics: {Loading:["Getting Topics..."]},
        }
      }

    // AXIOS request

    componentDidMount(){
        const url = `/plugins/topic_modeling`;
        let result = {
            messages: this.props.posts.map(a => a.message)
          };
        console.log(result)
        axios({
            url: url,
            method: 'post',
            data: result
            })
          .then((response) => {
            console.log(response.data.pets);
            this.setState({topics: response.data})
          })
          .catch(err => console.log(err));
        }    


    render() {
        
        return (
            <Accordion>
            {Object.keys(this.state.topics).map((visit, index) =>
            <AccordionItem>
                <AccordionItemHeading>
                    <AccordionItemButton>
                        {visit}  
                    </AccordionItemButton>
                </AccordionItemHeading>
                {this.state.topics[visit].map((text, i) =>
                <AccordionItemPanel>
                    <p>
                        {text}
                    </p>
                </AccordionItemPanel>
                )}
            </AccordionItem>
            )}
        </Accordion>
            
        );
    }
}

const style = {
    rhs: {
        padding: '10px',
    },
};


