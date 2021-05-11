from rest_framework import serializers
import status

from main.serializers.subtask.ser_subtask import SubtaskSerializer
from main.models import Subtask
from main.validation.val_auth import authenticate, authorization_error
from main.validation.val_custom import CustomAPIException


class UpdateSubtaskSerializer(SubtaskSerializer):
    id = serializers.IntegerField(error_messages={
        'invalid': 'Subtask ID must be a number.'
    })
    data = serializers.DictField(
        allow_empty=False,
        error_messages={
            'empty': 'Subtask data cannot be empty.'
        }
    )
    auth_user = serializers.CharField(allow_blank=True)
    auth_token = serializers.CharField(allow_blank=True)

    class Meta(SubtaskSerializer.Meta):
        fields = 'id', 'data', 'auth_user', 'auth_token'

    @staticmethod
    def validate_title(value):
        if not value:
            raise CustomAPIException('title',
                                     'Title cannot be empty.',
                                     status.HTTP_400_BAD_REQUEST)

    @staticmethod
    def validate_done(value):
        if value == '' or value is None or not value:
            raise CustomAPIException('done',
                                     'Done cannot be empty.',
                                     status.HTTP_400_BAD_REQUEST)

    @staticmethod
    def validate_order(value):
        if value == '' or value is None or not value:
            raise CustomAPIException('order',
                                     'Order cannot be empty.',
                                     status.HTTP_400_BAD_REQUEST)

    def validate(self, attrs):
        user = authenticate(attrs.get('auth_user'), attrs.get('auth_token'))

        subtask = Subtask.objects.select_related(
            'task',
            'task__user',
            'task__column__board'
        ).get(id=attrs.get('id'))

        if not user.is_admin and subtask.task.user != user \
                or subtask.task.column.board.team_id != user.team.id:
            raise authorization_error

        data = attrs.get('data')
        if 'title' in data.keys():
            self.validate_title(data.get('title'))
        if 'done' in data.keys():
            self.validate_done(data.get('done'))
        if 'order' in data.keys():
            self.validate_order(data.get('order'))

        serializer = SubtaskSerializer(subtask, data=data, partial=True)
        if serializer.is_valid(raise_exception=True):
            self.instance = subtask
            return data

    def to_representation(self, instance):
        return {'msg': 'Subtask update successful.',
                'id': instance.id}
