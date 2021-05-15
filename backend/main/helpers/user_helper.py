import bcrypt

from main.models import User


class UserHelper:
    def __init__(self, team):
        self.team = team
        self.__counter = 0

    def create_user(self, is_admin=False):
        """
        Creates a new user and returns a dictionary containing user data, as
        well as the token.
        """
        user = User.objects.create(
            username=(
                f'{"member" if is_admin else "admin"}-{self.__counter}'
            ),
            password=b'$2b$12$DKVJHUAQNZqIvoi.OMN6v.x1ZhscKhbzSxpOBMykHgTIMeeJ'
                     b'pC6me',
            is_admin=is_admin,
            team=self.team
        )
        self.__counter += 1
        return {'username': user.username,
                'password': user.password,
                'password_raw': 'barbarbar',
                'is_admin': user.is_admin,
                'team': user.team,
                'token': bcrypt.hashpw(
                    bytes(user.username, 'utf-8') + user.password,
                    bcrypt.gensalt()
                ).decode('utf-8')}


