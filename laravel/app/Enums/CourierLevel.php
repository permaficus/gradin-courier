<?php

namespace App\Enums;

enum CourierLevel: int
{
    case One = 1;
    case Two = 2;
    case Three = 3;
    case Four = 4;
    case Five = 5;

    public static function values(): array
    {
        return array_map(fn (self $level) => $level->value, self::cases());
    }
}
