from rest_framework import serializers
import status

from main.serializers.board.ser_board import BoardSerializer
from main.validation.val_auth import authenticate, authorize
from main.validation.val_custom import CustomAPIException
from main.models import Board


class UpdateBoardSerializer(serializers.ModelSerializer):
    id = serializers.IntegerField(error_messages={
        'null': 'Board ID cannot be null.',
        'invalid': 'Board ID must be a number.'
    })
    payload = serializers.DictField(allow_empty=False)
    auth_user = serializers.CharField(allow_blank=True)
    auth_token = serializers.CharField(allow_blank=True)

    class Meta:
        model = BoardSerializer.Meta.model
        fields = 'id', 'payload', 'auth_user', 'auth_token'

    def validate(self, attrs):
        user = authenticate(attrs.get('auth_user'), attrs.get('auth_token'))

        try:
            board = Board.objects.get(id=attrs.get('id'))
        except Board.DoesNotExist:
            raise CustomAPIException('board_id',
                                     'Board not found.',
                                     status.HTTP_404_NOT_FOUND)

        authorize(user, board.team_id)

        payload = attrs.get('payload')
        board_serializer = BoardSerializer(board, data=payload, partial=True)
        board_serializer.is_valid(raise_exception=True)

        self.instance = board
        return payload

    def to_representation(self, instance):
        return {'msg': 'Board updated successfuly.',
                'id': instance.id}
