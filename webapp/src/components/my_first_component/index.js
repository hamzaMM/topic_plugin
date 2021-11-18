import React from 'react';

export default class MyFirstComponent extends React.PureComponent {
    render() {
        const iconStyle = {
            display: 'inline-block',
            margin: '0 7px 0 1px',
        };
        const style = {
            margin: '.5em 0 .5em',
            padding: '0 12px 0 15px',
            backgroundColor: 'rgba(255,255,255,0.6)',
        };
        const url = 'https://developers.mattermost.com/extend/plugins/webapp/reference';
        return (
            <div style={style}>
                <i
                    className='icon fa fa-plug'
                    style={iconStyle}
                />
                <a href={url}>More info on plugins</a></div>
        );
    }
}