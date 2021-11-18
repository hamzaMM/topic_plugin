import React from 'react';

import RightSideBar from './components/right_hand_sidebar';
import MyFirstComponent from './components/my_first_component';

const Icon = () => <i className='icon fa fa-plug'/>;
class HelloWorldPlugin {
    initialize(registry, store) {
        registry.registerMainMenuAction(
            'Saying hi',
            () => allert('Hi again'),
        );

        registry.registerChannelHeaderButtonAction(
            <Icon/>,
            () => store.dispatch(toggleRHSPlugin),
            'Gets Topics',
        );

        const {toggleRHSPlugin} = registry.registerRightHandSidebarComponent(
            RightSideBar, 'Topics');

        const {leftsidebarheader} = registry.registerLeftSidebarHeaderComponent(MyFirstComponent);
    }
}

window.registerPlugin('topic_modeling', new HelloWorldPlugin());

