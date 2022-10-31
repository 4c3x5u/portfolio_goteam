from rest_framework import serializers

from server.main.models import Board


class BoardSerializer(serializers.ModelSerializer):
    class Meta:
        model = Board
        fields = ('id', 'team', 'name')
        extra_kwargs = {
            'name': {
                'error_messages': {
                    'blank': 'Board name cannot be blank.',
                    'null': 'Board name cannot be null.'
                }
            },
            'team': {
                'error_messages': {
                    'blank': 'Board team cannot be blank.',
                    'null': 'Board team cannot be null.',
                    'does_not_exist': 'Team does not exist.'
                }
            }
        }
