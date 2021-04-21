import React from 'react';
import { Col, Container, Row } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faGithub, faLinkedin } from '@fortawesome/free-brands-svg-icons';

import './footer.sass';

const Footer = () => (
  <div id="Footer">
    <Container>
      <Row className="Row">
        <Col className="Col" xs={6}>
          <a
            href="https://github.com/alicandev/portfolio_goteam"
            target="_blank"
            rel="noreferrer"
          >
            <FontAwesomeIcon icon={faGithub} />

            PROJECT
          </a>
        </Col>
        <Col className="Col" xs={6}>
          <a
            href="https://www.linkedin.com/in/4c3x5u/"
            target="_blank"
            rel="noreferrer"
          >
            <FontAwesomeIcon icon={faLinkedin} />

            AUTHOR
          </a>
        </Col>
      </Row>
    </Container>
  </div>
);

export default Footer;
