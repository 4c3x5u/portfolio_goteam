from rest_framework import serializers
import status

from main.serializers.board.base import BoardSerializer
from main.validation.auth import authenticate_custom, authorization_error, \
    authorize_custom
from main.validation.custom import CustomAPIException
from main.models import Board


class DeleteBoardSerializer(serializers.ModelSerializer):
    auth_user = serializers.CharField(allow_blank=True)
    auth_token = serializers.CharField(allow_blank=True)
    id = serializers.IntegerField(error_messages={
        'null': 'Board ID cannot be null.',
        'invalid': 'Board ID must be a number.'
    })

    class Meta:
        model = BoardSerializer.Meta.model
        fields = 'id', 'auth_user', 'auth_token',

    def validate(self, attrs):
        # authenticate
        auth_user = attrs.get('auth_user')
        auth_token = attrs.get('auth_token')
        user, authentication_error = authenticate_custom(auth_user, auth_token)
        if authentication_error:
            raise authentication_error

        # authorize
        authorize_error = authorize_custom(user.username)
        if authorize_error:
            raise authorize_error

        try:
            board = Board.objects.get(id=attrs.get('id'))
        except Board.DoesNotExist:
            raise CustomAPIException('board_id',
                                     'Board not found.',
                                     status.HTTP_404_NOT_FOUND)

        if board.team_id != user.team_id:
            raise authorization_error

        return board

    def delete(self):
        self.instance = {'id': self.validated_data.id}
        return self.validated_data.delete()

    def to_representation(self, instance):
        return {
            'msg': 'Board deleted successfully.',
            'id': instance.get('id'),
        }



