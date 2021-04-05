import React, { useState } from 'react';
import { Container, Row } from 'react-bootstrap';
import {
  faChalkboardTeacher, faQuestionCircle, faUsers,
} from '@fortawesome/free-solid-svg-icons';

import Controls from './Controls';

import './appheader.sass';
import logo from '../../../assets/homeHeader.svg';

const AppHeader = () => {
  const [teamControlsOn, setTeamControlsOn] = useState(false);
  const [boardControlsOn, setBoardControlsOn] = useState(false);

  const toggleTeamControls = () => {
    if (!teamControlsOn) { setBoardControlsOn(teamControlsOn); }
    setTeamControlsOn(!teamControlsOn);
  };

  const toggleBoardControls = () => {
    if (!boardControlsOn) { setTeamControlsOn(boardControlsOn); }
    setBoardControlsOn(!boardControlsOn);
  };

  return (
    <div id="AppHeader">
      <div className="Logo">
        <img alt="logo" src={logo} />
      </div>
      <div className="ControlsWrapper">
        <Container>
          <Row className="ControlsRow">
            <Controls
              name="team"
              toggle={toggleTeamControls}
              isToggled={teamControlsOn}
              icon={faUsers}
            />

            <Controls
              name="boards"
              toggle={toggleBoardControls}
              isToggled={boardControlsOn}
              icon={faChalkboardTeacher}
            />

            <Controls
              name="help"
              action={() => console.log('boards button clicked')}
              icon={faQuestionCircle}
            />
          </Row>
        </Container>
      </div>
    </div>
  );
};

export default AppHeader;
