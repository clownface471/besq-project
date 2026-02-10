// Centralized Role Configuration
export const ROLES = {
    MANAGER: 'MANAGER',         // Can access Layer 1, 2, 3
    TEAM_LEADER: 'LEADER',      // Can access Layer 2, 3
    OPERATOR: 'OPERATOR',       // Can access Layer 3 only
    ADMIN: 'ADMIN'              // Superuser
};

export const LABELS = {
    [ROLES.MANAGER]: 'Manager',
    [ROLES.TEAM_LEADER]: 'Team Leader',
    [ROLES.OPERATOR]: 'Operator'
};